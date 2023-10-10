package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var (
	urlMap     map[string]string
	serverAddr string
	baseURL    string
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		redirect(w, r)
	case http.MethodPost:
		handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

	r := chi.NewRouter()
	r.Get("/", mainPage)
	r.Post("/", handlePost)
	r.Get("/{id}", redirect)
	err := http.ListenAndServe(serverAddr, r)
	if err != nil {
		panic(err)
	}
}
