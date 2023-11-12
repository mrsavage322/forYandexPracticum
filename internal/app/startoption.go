package app

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var (
	ServerAddr   string
	BaseURL      string
	URLMap       URLStorage
	FilePATH     string
	DatabaseAddr string
	URLMapDB     URLDatabaseStorage
)

type Config struct {
	ServerAddr   string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
	FilePATH     string `env:"FILE_STORAGE_PATH"`
	DatabaseAddr string `env:"DATABASE_DSN"`
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
		if cfg.FilePATH != "" {
			FilePATH = cfg.FilePATH
		}
		if cfg.DatabaseAddr != "" {
			DatabaseAddr = cfg.DatabaseAddr
		}
	}
}

func SetFlags() {
	flag.StringVar(&ServerAddr, "a", "localhost:8080", "Address to run the HTTP server")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
	flag.StringVar(&FilePATH, "f", "/tmp/short-url-db.json", "Full path to the storage file")
	flag.StringVar(&DatabaseAddr, "d", "", "Address to connect with Database")
	flag.Parse()
}
