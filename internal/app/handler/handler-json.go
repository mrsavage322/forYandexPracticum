package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mrsavage322/foryandex/internal/app"
	"log"
	"net/http"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func HandleJSON(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req Request
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := req.URL
	id := GenerateRandomID(5)
	shortURL := fmt.Sprintf("%s/%s", app.Cfg.BaseURL, id)

	if app.Cfg.DatabaseAddr != "" {
		err := app.Cfg.URLMapDB.SetDB(context.Background(), id, link, app.Cfg.UserID)
		if err != nil {
			originalURL, err := app.Cfg.URLMapDB.GetReverse(context.Background(), link, app.Cfg.UserID)
			if err != nil {
				log.Println(err)
				return
			}
			shortURL := fmt.Sprintf("%s/%s", app.Cfg.BaseURL, originalURL)
			resp := Response{Result: shortURL}
			responseData, err := json.Marshal(resp)
			if err != nil {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write(responseData)
			return
		}
	} else {
		app.Cfg.URLMap.Set(id, link)
	}

	resp := Response{Result: shortURL}
	responseData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}
