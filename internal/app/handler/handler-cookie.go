package handler

import (
	"encoding/json"
	"net/http"
)

//	type ResponseBatchForUser struct {
//		OriginalURL string `json:"original_url"`
//		ShortURL    string `json:"short_url"`
//	}
const (
	cookieName = "user_id"
)

// GetUserURLs обрабатывает запрос на получение URL пользователя из базы данных
func GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userIDCookie, err := r.Cookie(cookieName)
	if err != nil || userIDCookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Здесь вы должны использовать ваш метод для получения URL из базы данных
	// Пример: urls, err := app.Cfg.URLMapDB.GetURLsByUserID(userIDCookie.Value)
	// Замените app.Cfg.URLMapDB.GetURLsByUserID(userIDCookie.Value) на ваш реальный метод получения данных из БД

	// Ваш код для запроса к базе данных и получения URL
	urls := []struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}{}

	// Проверяем, есть ли у пользователя сокращенные URL
	if len(urls) == 0 {
		// Отправляем статус 204 No Content, так как нет данных
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Отправляем данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}
