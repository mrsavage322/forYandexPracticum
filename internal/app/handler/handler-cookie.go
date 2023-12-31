package handler

import (
	"context"
	"encoding/json"
	"github.com/mrsavage322/foryandex/internal/app"
	"net/http"
)

const (
	cookieName = "user_id"
)

type ResponseBatchForUser struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func Authenticator(w http.ResponseWriter, r *http.Request) {
	userID, err := r.Cookie(cookieName)
	if err != nil || userID.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func GetUserURLs(w http.ResponseWriter, r *http.Request) {
	if app.Cfg.DatabaseAddr != "" {
		urlMap, err := app.Cfg.URLMapDB.GetDBAll(context.Background(), app.Cfg.UserID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else if len(urlMap) == 0 {
			http.Error(w, "Empty!", http.StatusNoContent)
			return
		}
		Authenticator(w, r)
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
