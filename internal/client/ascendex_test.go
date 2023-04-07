package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/yeahyeahcore/wws-task/internal/client"
	"github.com/yeahyeahcore/wws-task/internal/models"
)

func TestConnect(t *testing.T) {
	client := client.NewAscendex(context.Background(), &client.AscendexDeps{
		Logger:     logrus.New(),
		AuthKey:    "auth-key1",
		AuthSecret: "auth-signature1",
	})

	assert.NoError(t, client.Connection())
}

func TestClose(t *testing.T) {
	client := client.NewAscendex(context.Background(), &client.AscendexDeps{
		Logger:     logrus.New(),
		AuthKey:    "auth-key2",
		AuthSecret: "auth-signature2",
	})

	client.Connection()
	client.SubscribeToChannel("BTC_USDT")
	client.Disconnect()

	assert.Error(t, client.SubscribeToChannel("BTC_USDT"))
}

func TestSubscribeToChannel(t *testing.T) {
	client := client.NewAscendex(context.Background(), &client.AscendexDeps{
		Logger:     logrus.New(),
		AuthKey:    "auth-key3",
		AuthSecret: "auth-signature3",
	})

	client.Connection()

	assert.NoError(t, client.SubscribeToChannel("BTC_USDT"))
}

func TestReadMessagesFromChannel(t *testing.T) {
	client := client.NewAscendex(context.Background(), &client.AscendexDeps{
		Logger:     logrus.New(),
		AuthKey:    "BclE7dBGbS1AP3VnOuq6s8fJH0fWbH7r",
		AuthSecret: "fAZcQRUMxj3eX3DreIjFcPiJ9UR3ZTdgIw8mxddvtcDxLoXvdbXJuFQYadUUsF7q",
	})
	receiveCh := make(chan models.BestOrderBook)

	client.Connection()
	client.SubscribeToChannel("BTC_USDT")

	go client.ReadMessagesFromChannel(receiveCh)

	select {
	case bestOrderBook := <-receiveCh:
		assert.NotEmpty(t, bestOrderBook.Ask.Amount)
		assert.NotEmpty(t, bestOrderBook.Ask.Price)
		assert.NotEmpty(t, bestOrderBook.Bid.Amount)
		assert.NotEmpty(t, bestOrderBook.Bid.Price)
	case <-time.After(time.Second * 4):
		t.Errorf("timeout waiting for message")
		return
	}
}
