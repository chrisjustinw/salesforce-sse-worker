package main

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"salesforce-sse-worker/internal/di"
	"salesforce-sse-worker/internal/handler"
)

func main() {
	container, err := di.Provides()
	if err != nil {
		panic(err.Error())
	}

	e := echo.New()
	e.HideBanner = true

	if err := container.Invoke(func(messageHandler handler.ConversationHandler) {
		e.POST("/conversation/create", messageHandler.CreateConversation)
		e.PUT("/conversation/continuation-token", messageHandler.GenerateContinuationToken)

	}); err != nil {
		panic(err.Error())
	}

	if err := e.Start(":8888"); err != nil {
		slog.Error("Error handler", slog.String("err", err.Error()))
	}
}
