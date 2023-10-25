package app

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var sugarLogger *zap.SugaredLogger
var logger *zap.Logger

func InitializeLogger() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}
	defer logger.Sync()
	sugarLogger = logger.Sugar()
}

// LogRequest логирует информацию о запросе.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		responseWriter := NewStatusSizeLoggingResponseWriter(w)
		next.ServeHTTP(w, r)
		elapsed := time.Since(startTime)

		sugarLogger.Infow("Request",
			"URL", r.URL,
			"Status", responseWriter.Status(),
			"Method", r.Method,
			"Size", responseWriter.Size(),
			"Elapsed", elapsed,
		)
	})
}

func NewStatusSizeLoggingResponseWriter(w http.ResponseWriter) *StatusSizeLoggingResponseWriter {
	return &StatusSizeLoggingResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK, // По умолчанию, предполагаем успешный статус
		size:           0,
	}
}

type StatusSizeLoggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *StatusSizeLoggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *StatusSizeLoggingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *StatusSizeLoggingResponseWriter) Status() int {
	return w.status
}

func (w *StatusSizeLoggingResponseWriter) Size() int {
	return w.size
}
