package handler

import (
	"encoding/json"
	"fmt"
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
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		//// Асинхронное удаление URL из базы данных.
		//err = DeleteURLs(urls)
		//if err != nil {
		//	http.Error(w, "Error deleting URLs", http.StatusInternalServerError)
		//	return
		//}
		for _, url := range urls {
			err := app.Cfg.URLMapDB.DeleteDB(url, app.Cfg.UserID)
			if err != nil {
				log.Println("Problem to remove BD", err)
				return
			}
		}
		fmt.Println("URL:", urls)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("URLs marked as deleted"))
	}

	//// DeleteURLs выполняет асинхронное удаление URL из базы данных.
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

	//func MarkAsDeleted(id string) error {
	//	// Используем параметризованный SQL-запрос для безопасной вставки данных.
	//	_, err := db.Exec(context.Background(), "UPDATE url_storage SET is_deleted = $1 WHERE uuid = $2", true, id)
	//	if err != nil {
	//		return fmt.Errorf("error updating URL: %v", err)
	//	}
	//	return nil
	//}
}
