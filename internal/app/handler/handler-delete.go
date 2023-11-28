package handler

import (
	"encoding/json"
	"github.com/mrsavage322/foryandex/internal/app"
	"log"
	"net/http"
	"sync"
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
		var wg sync.WaitGroup

		go func() {
			defer wg.Done()
			for _, url := range urls {
				err := app.Cfg.URLMapDB.DeleteDBPrepare(url, app.Cfg.UserID)
				resultChan <- err
				if err != nil {
					log.Println("Problem to remove url from BD ", err)
					return
				}
			}
		}()
		go func() {
			wg.Wait()
			close(resultChan)
		}()
		for _, url := range urls {
			err := app.Cfg.URLMapDB.DeleteDBFinally(url, app.Cfg.UserID)
			resultChan <- err
			if err != nil {
				log.Println("Problem to remove url from BD ", err)
				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("URLs deleted"))
	}
}
