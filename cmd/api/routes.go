package main

import (
	_ "github.com/aerosystems/customer-service/docs" // docs are generated by Swag CLI, you have to import it.
	middleware "github.com/aerosystems/customer-service/internal/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (app *Config) NewRouter() *echo.Echo {
	e := echo.New()

	docsGroup := e.Group("/docs")
	docsGroup.Use(middleware.BasicAuthMiddleware)
	docsGroup.GET("/*", echoSwagger.WrapHandler)

	e.GET("/v1/customers", app.baseHandler.GetCustomer, middleware.AuthTokenMiddleware([]string{"customer", "support", "admin"}))

	return e
}
