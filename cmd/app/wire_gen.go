// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/customer-service/internal/config"
	"github.com/aerosystems/customer-service/internal/infrastructure/adapters/broker"
	"github.com/aerosystems/customer-service/internal/infrastructure/repository/fire"
	"github.com/aerosystems/customer-service/internal/presenters/http"
	"github.com/aerosystems/customer-service/internal/presenters/http/handlers"
	"github.com/aerosystems/customer-service/internal/usecases"
	"github.com/aerosystems/customer-service/pkg/logger"
	"github.com/aerosystems/customer-service/pkg/pubsub"
	"github.com/sirupsen/logrus"
)

// Injectors from wire.go:

//go:generate wire
func InitApp() *App {
	logger := ProvideLogger()
	logrusLogger := ProvideLogrusLogger(logger)
	config := ProvideConfig()
	baseHandler := ProvideBaseHandler(logrusLogger, config)
	client := ProvideFirestoreClient(config)
	customerRepo := ProvideFireCustomerRepo(client)
	pubSubClient := ProvidePubSubClient(config)
	subscriptionEventsAdapter := ProvideSubscriptionEventsAdapter(pubSubClient, config)
	customerUsecase := ProvideCustomerUsecase(logrusLogger, customerRepo, subscriptionEventsAdapter)
	customerHandler := ProvideCustomerHandler(logrusLogger, baseHandler, customerUsecase)
	server := ProvideHttpServer(logrusLogger, config, customerHandler)
	app := ProvideApp(logrusLogger, config, server)
	return app
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server) *App {
	app := NewApp(log, cfg, httpServer)
	return app
}

func ProvideLogger() *logger.Logger {
	loggerLogger := logger.NewLogger()
	return loggerLogger
}

func ProvideConfig() *config.Config {
	configConfig := config.NewConfig()
	return configConfig
}

func ProvideCustomerUsecase(log *logrus.Logger, customerRepo usecases.CustomerRepository, subscriptionEventsAdapter usecases.SubscriptionEventsAdapter) *usecases.CustomerUsecase {
	customerUsecase := usecases.NewCustomerUsecase(log, customerRepo, subscriptionEventsAdapter)
	return customerUsecase
}

func ProvideFireCustomerRepo(client *firestore.Client) *fire.CustomerRepo {
	customerRepo := fire.NewCustomerRepo(client)
	return customerRepo
}

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, customerHandler *handlers.CustomerHandler) *HttpServer.Server {
	server := HttpServer.NewServer(log, customerHandler)
	return server
}

func ProvideCustomerHandler(log *logrus.Logger, baseHandler *handlers.BaseHandler, customerUsecase handlers.CustomerUsecase) *handlers.CustomerHandler {
	customerHandler := handlers.NewCustomerHandler(baseHandler, customerUsecase)
	return customerHandler
}

// wire.go:

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideSubscriptionEventsAdapter(pubSubClient *PubSub.Client, cfg *config.Config) *broker.SubscriptionEventsAdapter {
	return broker.NewSubscriptionEventsAdapter(pubSubClient, cfg.SubscriptionTopicId, cfg.SubscriptionSubName, cfg.SubscriptionCreateFreeTrialEndpoint, cfg.SubscriptionServiceApiKey)
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvidePubSubClient(cfg *config.Config) *PubSub.Client {
	client, err := PubSub.NewClientWithAuth(cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *handlers.BaseHandler {
	return handlers.NewBaseHandler(log, cfg.Mode)
}
