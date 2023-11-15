package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if Cfg.DatabaseAddr != "" {
		originalURL, err := Cfg.URLMapDB.Get(id)
		if err != nil {
			http.Error(w, "Non-existent identifier", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		originalURL, err := Cfg.URLMap.Get(id)
		if err != nil {
			http.Error(w, "Non-existent identifier", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
