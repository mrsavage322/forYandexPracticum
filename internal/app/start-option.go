package app

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/mrsavage322/foryandex/internal/storage"
)

var (
	ServerAddr string
	BaseURL    string
	URLMap     storage.URLStorage
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func SetConfig() {
	var cfg Config
	er := env.Parse(&cfg)
	if er == nil {
		if cfg.ServerAddr != "" {
			ServerAddr = cfg.ServerAddr
		}
		if cfg.BaseURL != "" {
			BaseURL = cfg.BaseURL
		}
	}
}

func SetFlags() {
	flag.StringVar(&ServerAddr, "a", "localhost:8080", "Address to run the HTTP server")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
	flag.Parse()
}
