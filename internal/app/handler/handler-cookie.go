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

		var response []ResponseBatchForUser
		for shortURL, originalURL := range urlMap {
			resp := ResponseBatchForUser{
				OriginalURL: originalURL,
				ShortURL:    app.Cfg.BaseURL + "/" + shortURL,
			}
			response = append(response, resp)
		}

		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	}
}
