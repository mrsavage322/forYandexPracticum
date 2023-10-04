package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urlMap map[string]string

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод должен быть POST", http.StatusBadRequest)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

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
	urlMap[id] = link

	// Возвращаем сокращенный URL как ответ с кодом 201
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, ok := urlMap[id]
	if !ok {
		http.Error(w, "Несуществующий идентификатор", http.StatusBadRequest)
		return
	}

	// Выполняем перенаправление на оригинальный URL с кодом 307
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateRandomID(length int) string {
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func main() {
	urlMap = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", mainPage)
	mux.HandleFunc("/", redirect)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
