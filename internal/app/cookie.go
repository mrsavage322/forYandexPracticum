package app

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GiveCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isGzipSupported(r) {
		}

		next.ServeHTTP(w, r)
	})
}
