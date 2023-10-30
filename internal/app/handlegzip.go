package app

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (grw GzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

func isGzipSupported(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isGzipSupported(r) {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			next.ServeHTTP(GzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
