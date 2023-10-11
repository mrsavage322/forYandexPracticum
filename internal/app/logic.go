package app

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var (
	UrlMap     map[string]string
	ServerAddr string
	BaseURL    string
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	link := strings.TrimSpace(string(bodyBytes))
	if link == "" {
		http.Error(w, "Empty link", http.StatusBadRequest)
		return
	}

	id := generateRandomID(5)
	shortURL := fmt.Sprintf("%s/%s", BaseURL, id)
	UrlMap[id] = link

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalURL, ok := UrlMap[id]
	if !ok {
		http.Error(w, "Non-existent identifier", http.StatusBadRequest)
		return
	}

	// Выполняем перенаправление на оригинальный URL с кодом 307
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateRandomID(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
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
