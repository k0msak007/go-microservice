package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/k0msak007/go-microservice/src/config"
	httpcontrollers "github.com/k0msak007/go-microservice/src/controllers/httpControllers"
	"github.com/k0msak007/go-microservice/src/repositories"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type IHttpServer interface {
	StartUserServer()
	StartItemServer()
}

type httpServer struct {
	app      *echo.Echo
	cfg      *config.Config
	dbClient *mongo.Client
}

func NewHttpServer(cfg *config.Config, dbClient *mongo.Client) IHttpServer {
	return &httpServer{
		app:      echo.New(),
		cfg:      cfg,
		dbClient: dbClient,
	}
}

func (s *httpServer) listener() {
	log.Printf("server is starting on %s", s.cfg.App.Url)
	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		s.app.Logger.Fatal("shutting down the server")
	}
}

func (s *httpServer) gracefullyShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.app.Shutdown(ctx); err != nil {
		s.app.Logger.Fatal(err)
	}
}

func (s *httpServer) StartUserServer() {
	usersController := httpcontrollers.UserHttpController{
		Cfg: s.cfg,
		UserRepository: &repositories.UserRepository{
			Client: s.dbClient,
		},
	}

	s.app.GET("/user/:user_id", usersController.FindOneUser)

	// Server
	go s.listener()
	s.gracefullyShutdown()

}

func (s *httpServer) StartItemServer() {
	itemController := httpcontrollers.ItemHttpController{
		ItemRepository: &repositories.ItemRepository{
			Client: s.dbClient,
		},
	}

	s.app.GET("/item", itemController.FindItems)
	s.app.GET("/item/:item_id", itemController.FindOneItem)

	// Server
	go s.listener()
	s.gracefullyShutdown()
}
