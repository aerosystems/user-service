//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/customer-service/internal/adapters"
	"github.com/aerosystems/customer-service/internal/common/config"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	HttpServer "github.com/aerosystems/customer-service/internal/presenters/http"
	"github.com/aerosystems/customer-service/internal/presenters/http/handlers"
	"github.com/aerosystems/customer-service/internal/usecases"
	"github.com/aerosystems/customer-service/pkg/logger"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(handlers.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*adapters.FirestoreCustomerRepo)),
		wire.Bind(new(usecases.SubscriptionAdapter), new(*adapters.SubscriptionAdapter)),
		wire.Bind(new(usecases.ProjectAdapter), new(*adapters.ProjectAdapter)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideLogrusLogger,
		ProvideFirestoreClient,
		ProvideCustomerUsecase,
		ProvideFirestoreCustomerRepo,
		ProvideHttpServer,
		ProvideCustomerHandler,
		ProvideEchoErrorHandler,
		ProvideSubscriptionAdapter,
		ProvideProjectAdapter,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideSubscriptionAdapter(cfg *config.Config) *adapters.SubscriptionAdapter {
	subscriptionAdapter, err := adapters.NewSubscriptionAdapter(cfg.ProjectServiceGRPCAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return subscriptionAdapter
}

func ProvideProjectAdapter(cfg *config.Config) *adapters.ProjectAdapter {
	projectAdapter, err := adapters.NewProjectAdapter(cfg.ProjectServiceGRPCAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return projectAdapter
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideCustomerUsecase(log *logrus.Logger, customerRepo usecases.CustomerRepository, subscriptionAdapter usecases.SubscriptionAdapter, projectAdapter usecases.ProjectAdapter) *usecases.CustomerUsecase {
	panic(wire.Build(usecases.NewCustomerUsecase))
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirestoreCustomerRepo(client *firestore.Client) *adapters.FirestoreCustomerRepo {
	panic(wire.Build(adapters.NewFirestoreCustomerRepo))
}

func ProvideHttpServer(cfg *config.Config, log *logrus.Logger, customErrorHandler *echo.HTTPErrorHandler, customerHandler *handlers.FirebaseHandler) *HttpServer.Server {
	return HttpServer.NewServer(cfg.Port, log, customErrorHandler, customerHandler)
}

func ProvideCustomerHandler(log *logrus.Logger, customerUsecase handlers.CustomerUsecase) *handlers.FirebaseHandler {
	panic(wire.Build(handlers.NewFirebaseHandler))
}

func ProvideEchoErrorHandler(cfg *config.Config) *echo.HTTPErrorHandler {
	customErrorHandler := CustomErrors.NewEchoErrorHandler(cfg.Mode)
	return &customErrorHandler
}
