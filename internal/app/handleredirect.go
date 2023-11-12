package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalURL, ok := URLMapDB.Get(id)
	if ok == nil {
		http.Error(w, "Non-existent identifier", http.StatusBadRequest)
		return
	}

	// Выполняем перенаправление на оригинальный URL с кодом 307
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
