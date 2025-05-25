package main

import (
	"context"
	"salesforce-sse-worker/internal/di"
	"salesforce-sse-worker/internal/kafka"
)

func main() {
	container, err := di.Provides()
	if err != nil {
		panic(err.Error())
	}

	if err := container.Invoke(func(consumer kafka.Consumer) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		consumer.Consume(ctx)

		return
	}); err != nil {
		panic(err.Error())
	}
}
