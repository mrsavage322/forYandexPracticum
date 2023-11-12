package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if DatabaseAddr != "" {
		originalURL, ok := URLMapDB.Get(id)
		if ok == nil {
			http.Error(w, "Non-existent identifier", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		originalURL, ok := URLMap.Get(id)
		if !ok {
			http.Error(w, "Non-existent identifier", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
