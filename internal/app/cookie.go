package app

import (
	"github.com/google/uuid"
	"net/http"
)

const (
	cookieName     = "user_id"
	cookieHttpOnly = true
)

// AuthMiddleware is a middleware for authenticating users and setting a signed cookie with a unique user ID.
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
				HttpOnly: cookieHttpOnly,
			}

			http.SetCookie(w, &cookie)
		}

		next.ServeHTTP(w, r)
	})
}

// AuthenticatorMiddleware is a middleware for checking the authenticity of the user ID in the cookie.
func AuthenticatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie(cookieName)
		if err != nil || userID.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Your additional authentication logic goes here if needed
		Cfg.UserID = userID.Value
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Wrap your existing router with the new middleware
//func ActivateCookie(r *chi.Mux) http.Handler {
//	r.Use(AuthMiddleware)
//	r.Use(AuthenticatorMiddleware)
//	return r
//}
