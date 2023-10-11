package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var (
	urlMap     map[string]string
	serverAddr string
	baseURL    string
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func handlePost(w http.ResponseWriter, r *http.Request) {
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
	shortURL := fmt.Sprintf("%s/%s", baseURL, id)
	urlMap[id] = link

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalURL, ok := urlMap[id]
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

func main() {
	urlMap = make(map[string]string)

	flag.StringVar(&serverAddr, "a", "localhost:8080", "Address to run the HTTP server")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Base URL for shortened links")
	flag.Parse()

	var cfg Config
	er := env.Parse(&cfg)
	if er != nil {
		log.Print(er)
	} else {
		serverAddr = cfg.ServerAddr
		baseURL = cfg.BaseURL
	}

	r := chi.NewRouter()
	r.Get("/", redirect)
	r.Get("/{id}", redirect)
	r.Post("/", handlePost)
	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		panic(err)
	}
}
