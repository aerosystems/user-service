package HttpServer

import (
	"fmt"
	"github.com/aerosystems/customer-service/internal/presenters/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	port            int
	log             *logrus.Logger
	echo            *echo.Echo
	firebaseHandler *handlers.FirebaseHandler
}

func NewServer(
	port int,
	log *logrus.Logger,
	errorHandler *echo.HTTPErrorHandler,
	customerHandler *handlers.FirebaseHandler,

) *Server {
	server := &Server{
		port:            port,
		log:             log,
		echo:            echo.New(),
		firebaseHandler: customerHandler,
	}
	if errorHandler != nil {
		server.echo.HTTPErrorHandler = *errorHandler
	}
	return server
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	return s.echo.Start(fmt.Sprintf(":%d", s.port))
}
