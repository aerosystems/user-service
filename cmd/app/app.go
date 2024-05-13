package main

import (
	"github.com/aerosystems/customer-service/internal/config"
	"github.com/aerosystems/customer-service/internal/presenters/consumer"
	HttpServer "github.com/aerosystems/customer-service/internal/presenters/http"
	"github.com/sirupsen/logrus"
)

type App struct {
	log          *logrus.Logger
	cfg          *config.Config
	httpServer   *HttpServer.Server
	authConsumer *consumer.AuthSubscription
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
	httpServer *HttpServer.Server,
	authConsumer *consumer.AuthSubscription,
) *App {
	return &App{
		log:          log,
		cfg:          cfg,
		httpServer:   httpServer,
		authConsumer: authConsumer,
	}
}
