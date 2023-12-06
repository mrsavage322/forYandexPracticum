package app

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var Cfg Config

type Config struct {
	ServerAddr   string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
	FilePATH     string `env:"FILE_STORAGE_PATH"`
	DatabaseAddr string `env:"DATABASE_DSN"`
	URLMap       URLStorage
	URLMapDB     URLStorage
	UserID       string
}

func SetConfig() {
	err := env.Parse(&Cfg)
	if err != nil {
		return
	}
}

func SetFlags() {
	flag.StringVar(&Cfg.ServerAddr, "a", "localhost:8080", "Address to run the HTTP server")
	flag.StringVar(&Cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
	flag.StringVar(&Cfg.FilePATH, "f", "/tmp/short-url-db.json", "Full path to the storage file")
	flag.StringVar(&Cfg.DatabaseAddr, "d", "", "Address to connect with Database")
	flag.Parse()
}
