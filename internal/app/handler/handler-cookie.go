package handler

import (
	"encoding/json"
	"fmt"
	"github.com/mrsavage322/foryandex/internal/app"
	"net/http"
)

type ResponseBatchForUser struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func HandleURLsToUser(w http.ResponseWriter, r *http.Request) {
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
			app.Cfg.URLMapDB.Get(id)
		} else {
			app.Cfg.URLMap.Get(id)
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
