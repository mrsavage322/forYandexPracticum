package main

import (
	"net/http"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`https://practicum.yandex.ru/`))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(`http://localhost:8080/EwHXdJfB `))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, GetHandler)
	mux.HandleFunc(`/EwHXdJfB`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
