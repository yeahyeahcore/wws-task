package app

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/yeahyeahcore/wws-task/internal/client"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	ascendexClient := client.NewAscendex(ctx, &client.AscendexDeps{
		Logger:     logrus.New(),
		AuthKey:    "BclE7dBGbS1AP3VnOuq6s8fJH0fWbH7r",
		AuthSecret: "fAZcQRUMxj3eX3DreIjFcPiJ9UR3ZTdgIw8mxddvtcDxLoXvdbXJuFQYadUUsF7q",
	})

	defer cancel()

	ascendexClient.Connection()
}
