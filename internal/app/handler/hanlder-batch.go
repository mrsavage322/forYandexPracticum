package handler

import (
	"encoding/json"
	"fmt"
	"github.com/mrsavage322/foryandex/internal/app"
	"net/http"
)

type RequestBatch struct {
	CorrelID  string `json:"correlation_id"`
	OriginURL string `json:"original_url"`
}

type ResponseBatch struct {
	CorrelID string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

func HandleBatch(w http.ResponseWriter, r *http.Request) {
	var reqs []RequestBatch
	var resps []ResponseBatch

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, req := range reqs {
		link := req.OriginURL
		id := GenerateRandomID(5)
		shortURL := fmt.Sprintf("%s/%s", app.Cfg.BaseURL, id)
		correlationID := req.CorrelID

		if app.Cfg.DatabaseAddr != "" {
			app.Cfg.URLMapDB.Set(id, link)
		} else {
			app.Cfg.URLMap.Set(id, link)
		}

		resp := ResponseBatch{
			CorrelID: correlationID,
			ShortURL: shortURL,
		}
		resps = append(resps, resp)
	}

	responseData, err := json.Marshal(resps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}
