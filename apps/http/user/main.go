package main

import (
	"context"
	"os"

	"github.com/k0msak007/go-microservice/pkg/db"
	"github.com/k0msak007/go-microservice/src/config"
	"github.com/k0msak007/go-microservice/src/server"
)

func main() {
	cfg := config.NewConfig(func() string {
		if len(os.Args) > 1 {
			return os.Args[1]
		}
		return "./.env.http.user"
	}())

	dbClient := db.DbConn(cfg)
	defer dbClient.Disconnect(context.Background())

	server.NewHttpServer(cfg, dbClient).StartUserServer()
}
