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
		ok := URLMapDB.Set(id, link)
		if !ok {
			originalURL, okk := URLMapDB.GetReverse(link)
			if !okk {
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

//func (s *URLDBStorage) Get(key string) (string, bool) {
//	var originalURL string
//	err := s.conn.QueryRow(context.Background(), "SELECT short_url FROM url_storage WHERE original_url = $1", key).Scan(&originalURL)
//	if err != nil {
//		return "", false
//	}
//	return originalURL, true
//}
