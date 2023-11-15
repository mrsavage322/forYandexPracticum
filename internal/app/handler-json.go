package app

import (
	"encoding/json"
	"fmt"
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
	shortURL := fmt.Sprintf("%s/%s", BaseURL, id)

	if DatabaseAddr != "" {
		err := URLMapDB.Set(id, link)
		if err != nil {
			originalURL, err := URLMapDB.GetReverse(link)
			if err != nil {
				sugar.Warnln(err)
				return
			}
			shortURL := fmt.Sprintf("%s/%s", BaseURL, originalURL)
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
		URLMap.Set(id, link)
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
