package db

import (
	"context"
	"log"
	"time"

	"github.com/k0msak007/go-microservice/src/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DbConn(cfg *config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Db.Url))
	if err != nil {
		log.Fatalf("connect to mongodb url: %s failed: %v", cfg.Db.Url, err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("ping to mongodb url: %s failed: %v", cfg.Db.Url, err)
	}

	log.Printf("connected to mongodb url: %s", cfg.Db.Url)
	return client
}
