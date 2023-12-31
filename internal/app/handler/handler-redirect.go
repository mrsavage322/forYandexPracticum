package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrsavage322/foryandex/internal/app"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if app.Cfg.DatabaseAddr != "" {
		originalURL, err := app.Cfg.URLMapDB.GetDBNoCookie(id)
		if err != nil {
			//TODO Нужно поправить лоигку обработки с ссылкой, которой не было в БД и с ссылкой, которая была удалена
			http.Error(w, "Non-existent identifier", http.StatusGone)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		originalURL, err := app.Cfg.URLMap.Get(id)
		if err != nil {
			http.Error(w, "Non-existent identifier", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
