package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

var (
	urlMap     = make(map[string]string)
	urlMapLock sync.Mutex
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод должен быть POST", http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := strings.TrimSpace(string(bodyBytes))
	if link == "" {
		http.Error(w, "Пустая ссылка", http.StatusBadRequest)
		return
	}

	// Генерируем пять случайных букв для идентификатора
	id := generateRandomID(5)

	// Формируем сокращенный URL
	shortURL := fmt.Sprintf("http://localhost:8080/%s", id)

	// Сохраняем сокращенный URL в карту
	urlMapLock.Lock()
	urlMap[id] = link
	urlMapLock.Unlock()

	// Возвращаем сокращенный URL как ответ с кодом 201
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	urlMapLock.Lock()
	originalURL, ok := urlMap[id]
	urlMapLock.Unlock()
	if !ok {
		http.Error(w, "Несуществующий идентификатор", http.StatusBadRequest)
		return
	}

	// Выполняем перенаправление на оригинальный URL с кодом 307
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/shorten/", redirect)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
