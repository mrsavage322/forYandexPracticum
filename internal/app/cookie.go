package app

import (
	"github.com/google/uuid"
	"net/http"
)

const (
	cookieName     = "user_id"
	cookieHTTPOnly = true
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie(cookieName)
		if err != nil || userID.Value == "" {
			// Generate a new user ID
			newUserID := uuid.New().String()

			// Create a new cookie with the user ID
			cookie := http.Cookie{
				Name:     cookieName,
				Value:    newUserID,
				HttpOnly: cookieHTTPOnly,
			}

			http.SetCookie(w, &cookie)

			Cfg.UserID = newUserID
		} else {
			//Cfg.UserID = userID.Value
		}

		next.ServeHTTP(w, r)
	})
}

func AuthenticatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie(cookieName)
		if err != nil || userID.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
