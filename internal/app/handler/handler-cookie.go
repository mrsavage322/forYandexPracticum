package handler

import (
	"encoding/json"
	"github.com/mrsavage322/foryandex/internal/app"
	"net/http"
)

type ResponseBatchForUser struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// GetUserURLs обрабатывает запрос на получение URL пользователя из базы данных
func GetUserURLs(w http.ResponseWriter, r *http.Request) {
	if app.Cfg.DatabaseAddr != "" {
		urlMap, err := app.Cfg.URLMapDB.GetDBAll(app.Cfg.UserID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if len(urlMap) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(urlMap); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
