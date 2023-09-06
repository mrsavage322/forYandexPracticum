package main

import (
	"io/ioutil"
	"net/http"
)

var urlMap = make(map[string]string)

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Чтение URL из тела POST-запроса.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read URL from request body", http.StatusInternalServerError)
		return
	}
	url := string(body)

	// Генерация уникального идентификатора (короткого пути) и сохранение URL.
	id := generateID()
	urlMap[id] = url

	// Отправка ответа с коротким URL.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(id))
}

func generateID() string {
	// Здесь может быть логика для генерации уникального ID,
	// например, с использованием случайной генерации или хеширования.
	// В данном примере, используется фиксированный ID.
	return "EwHXdJfB"
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение идентификатора из URL.
	id := r.URL.Path[1:] // Удаление первого символа "/", чтобы получить идентификатор.

	// Поиск URL по идентификатору в карту urlMap.
	url, exists := urlMap[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Отправка ответа с перенаправлением и оригинальным URL в заголовке Location.
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc(`/EwHXdJfB`, GetHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
