package handler

import (
	"encoding/json"
	"github.com/mrsavage322/foryandex/internal/app"
	"log"
	"net/http"
)

// DeleteURLsHandler обрабатывает запрос на удаление URL.
func DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string
	if app.Cfg.DatabaseAddr != "" {

		// Декодирование JSON из тела запроса.
		err := json.NewDecoder(r.Body).Decode(&urls)
		if err != nil {
			//TODO Нужно поправить лоигку обработки с ссылкой, которой не было в БД и с ссылкой, которая была удалена
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		//// Асинхронное удаление URL из базы данных.
		//err = DeleteURLs(urls)
		//if err != nil {
		//	http.Error(w, "Error deleting URLs", http.StatusInternalServerError)
		//	return
		//}
		go func() {
			for _, url := range urls {
				err := app.Cfg.URLMapDB.DeleteDB(url, app.Cfg.UserID)
				if err != nil {
					log.Println("Problem to remove BD", err)
					return
				}
			}
		}()

	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("URLs deleted"))
}

//func DeleteURLs(ids []string) error {
//	// Выполнение асинхронного удаления URL в горутине.
//	go func() {
//		for _, id := range ids {
//			err := MarkURLAsDeleted(id)
//			if err != nil {
//				log.Printf("Error marking URL as deleted: %v\n", err)
//			}
//		}
//	}()
//
//	return nil
//}
