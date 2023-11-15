package app

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	link := strings.TrimSpace(string(bodyBytes))
	if link == "" {
		http.Error(w, "Empty link", http.StatusBadRequest)
		return
	}

	id := GenerateRandomID(5)
	shortURL := fmt.Sprintf("%s/%s", BaseURL, id)

	if DatabaseAddr != "" {
		err := URLMapDB.Set(id, link)
		if err != nil {
			originalURL, err := URLMapDB.GetReverse(link)
			if err != nil {
				return
			}
			shortURL := fmt.Sprintf("%s/%s", BaseURL, originalURL)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortURL))
			return
		}
	} else {
		URLMap.Set(id, link)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func GenerateRandomID(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
