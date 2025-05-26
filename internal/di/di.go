package di

import (
	"errors"
	"go.uber.org/dig"
	"salesforce-sse-worker/configs"
	"salesforce-sse-worker/internal/handler"
	"salesforce-sse-worker/internal/library"
	"salesforce-sse-worker/internal/repository"
	"salesforce-sse-worker/internal/service"
	"salesforce-sse-worker/internal/service/outbound"
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
	r.provide(configs.NewMongoClientConfig)
	r.provide(configs.NewSalesforceConfig)

	r.provide(library.NewHTTPClient)
	r.provide(library.NewKafkaProducer)
	r.provide(library.NewKafkaConsumer)
	r.provide(library.NewMongoDatabase)

	r.provide(repository.NewConversationMappingRepository)

	r.provide(handler.NewKafkaHandler)
	r.provide(handler.NewConversationHandler)

	r.provide(outbound.NewSalesforceOutbound)

	r.provide(service.NewConversationService)
}
