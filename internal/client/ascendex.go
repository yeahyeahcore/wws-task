package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/yeahyeahcore/wws-task/internal/models"
	"github.com/yeahyeahcore/wws-task/internal/utils"
	"github.com/yeahyeahcore/wws-task/pkg/json"
)

type APIClient interface {
	/*
		Implement a websocket connection function
	*/
	Connection() error

	/*
		Implement a disconnect function from websocket
	*/
	Disconnect()

	/*
		Implement a function that will subscribe to updates
		of BBO for a given symbol

		The symbol must be of the form "TOKEN_ASSET"
		As an example "USDT_BTC" where USDT is TOKEN and BTC is ASSET

		You will need to convert the symbol in such a way that
		it complies with the exchange standard
	*/
	SubscribeToChannel(symbol string) error

	/*
		Implement a function that will write the data that
		we receive from the exchange websocket to the channel
	*/
	ReadMessagesFromChannel(ch chan<- models.BestOrderBook)

	/*
		Implement a function that will support connecting to a websocket
	*/
	WriteMessagesToChannel()
}

type AscendexDeps struct {
	Logger     *logrus.Logger
	AuthKey    string
	AuthSecret string
}

type Ascendex struct {
	logger              *logrus.Logger
	ctx                 context.Context
	websocketConnection *websocket.Conn
	authKey             string
	authSecret          string
}

func NewAscendex(ctx context.Context, deps *AscendexDeps) APIClient {
	return &Ascendex{
		ctx:        ctx,
		logger:     deps.Logger,
		authKey:    deps.AuthKey,
		authSecret: deps.AuthSecret,
	}
}

func (receiver *Ascendex) Connection() error {
	streamURL := url.URL{Scheme: "wss", Host: "ascendex.com", Path: "/api/pro/v1/stream"}
	signature, timestamp := utils.GenerateSignature(receiver.authSecret)

	headers := http.Header{
		"x-auth-key":       []string{receiver.authKey},
		"x-auth-timestamp": []string{fmt.Sprint(timestamp)},
		"x-auth-signature": []string{signature},
	}

	websocketConnection, _, err := websocket.DefaultDialer.Dial(streamURL.String(), headers)
	if err != nil {
		return err
	}

	receiver.websocketConnection = websocketConnection

	return nil
}

func (receiver *Ascendex) Disconnect() {
	receiver.logger.Infoln("Ascendex API discounnecting...")
	receiver.websocketConnection.Close()
}

func (receiver *Ascendex) SubscribeToChannel(symbol string) error {
	json := map[string]interface{}{
		"op": "sub",
		"ch": []string{fmt.Sprintf("bbo:%s", symbol)},
	}

	if err := receiver.websocketConnection.WriteJSON(json); err != nil {
		receiver.logger.Errorf("Failed to write json on <SubscribeToChannel>: %v", err)
		return err
	}

	return nil
}

func (receiver *Ascendex) ReadMessagesFromChannel(ch chan<- models.BestOrderBook) {
	readFailedCount := 0

	defer close(ch)

	for {
		_, message, err := receiver.websocketConnection.ReadMessage()
		if err != nil {
			receiver.logger.Errorf("Failed to read message (%d) %v", readFailedCount, err)
			readFailedCount++
			continue
		}

		if readFailedCount > 5 {
			receiver.logger.Errorln("Reading message failed... leaving...")
			return
		}

		readFailedCount = 0
		reader := bytes.NewReader(message)

		bestOrderBook, err := json.Parse[models.BestOrderBook](reader)
		if err != nil {
			continue
		}

		ch <- *bestOrderBook
	}
}

func (receiver *Ascendex) WriteMessagesToChannel() {
	ticker := time.NewTicker(5 * time.Second)
	failedCount := 0

	defer ticker.Stop()

	for {
		select {
		case <-receiver.ctx.Done():
			return
		case <-ticker.C:
			if err := receiver.websocketConnection.WriteMessage(websocket.PingMessage, []byte(time.Now().String())); err != nil {
				receiver.logger.Errorf("Error sending ping message (%d): %v", failedCount, err)
				failedCount++
			}
			if failedCount > 5 {
				receiver.logger.Errorln("Sending ping message failed... leaving...")
				return
			}

			failedCount = 0
		}
	}
}
