package di

import (
	"errors"
	"go.uber.org/dig"
	"salesforce-sse-worker/configs"
	"salesforce-sse-worker/internal/handler"
	"salesforce-sse-worker/internal/kafka"
	"salesforce-sse-worker/internal/outbound"
	"salesforce-sse-worker/internal/outbound/client"
	"salesforce-sse-worker/internal/repository"
	"salesforce-sse-worker/internal/service"
)

type registry struct {
	container *dig.Container
	errors    []error
}

func (r *registry) provide(fn interface{}, opts ...dig.ProvideOption) {
	err := r.container.Provide(fn, opts...)
	r.errors = append(r.errors, err)
}

func (r *registry) GetError() error {
	return errors.Join(r.errors...)
}

func Provides() (*dig.Container, error) {
	r := registry{container: dig.New(), errors: []error{}}
	provides(&r)

	if err := r.GetError(); err != nil {
		return nil, err
	}

	return r.container, nil
}

func provides(r *registry) {
	r.provide(configs.ReadEnvFile)
	r.provide(configs.NewKafkaConfig)
	r.provide(configs.NewSaramaConfig)
	r.provide(configs.NewMongoConfig)
	r.provide(configs.NewMongoClient)
	r.provide(configs.NewMongoDatabase)
	r.provide(configs.NewSalesforceConfig)

	r.provide(repository.NewBaseRepository)
	r.provide(repository.NewConversationMappingRepository)

	r.provide(kafka.NewProducer)
	r.provide(kafka.NewConsumer)

	r.provide(handler.NewKafkaHandler)
	r.provide(handler.NewConversationHandler)

	r.provide(client.NewHTTPClient)
	r.provide(outbound.NewSalesforceOutbound)

	r.provide(service.NewConversationService)
}
