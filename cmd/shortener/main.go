package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrsavage322/foryandex/internal/app"
	"net/http"
)

func main() {
	app.URLMap = make(map[string]string)
	app.SetFlags()
	app.SetConfig()

	r := chi.NewRouter()
	r.Get("/", app.Redirect)
	r.Get("/{id}", app.Redirect)
	r.Post("/", app.HandlePost)

	err := http.ListenAndServe(app.ServerAddr, r)
	if err != nil {
		panic(err)
	}
}
