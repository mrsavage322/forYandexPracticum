package handler

import (
	"encoding/json"
	"github.com/mrsavage322/foryandex/internal/app"
	"log"
	"net/http"
)

func DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string
	if app.Cfg.DatabaseAddr != "" {

		err := json.NewDecoder(r.Body).Decode(&urls)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		resultChan := make(chan error, len(urls))

		go func() {
			for _, url := range urls {
				err := app.Cfg.URLMapDB.DeleteDB(url, app.Cfg.UserID)
				resultChan <- err
				if err != nil {
					log.Println("Problem to remove url from BD ", err)
					return
				}
			}
			close(resultChan)
		}()

	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("URLs deleted"))
}
