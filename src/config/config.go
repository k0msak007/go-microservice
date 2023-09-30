package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App  *App
	Db   *Db
	Grpc *Grpc
}

func NewConfig(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		App: &App{
			Url:     os.Getenv("APP_URL"),
			Appname: os.Getenv("APP_NAME"),
		},
		Db: &Db{
			Url: os.Getenv("DB_URL"),
		},
		Grpc: &Grpc{
			ItemAppUrl: os.Getenv("GRPC_ITEM_APP_URL"),
		},
	}
}

type App struct {
	Url     string
	Appname string
}

type Db struct {
	Url string
}

type Grpc struct {
	ItemAppUrl string
}
