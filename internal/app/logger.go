package app

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var sugar zap.SugaredLogger

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func InitializeLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Can not initialize zap logger: " + err.Error())
	}
	defer logger.Sync()

	sugar = *logger.Sugar()
}

func LogRequest(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		uri := r.RequestURI

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(startTime)

		sugar.Infow("Request",
			"URI", uri,
			"Method", r.Method,
			"Duration", duration,
		)

		sugar.Infow("Response",
			"Status", responseData.status,
			"Size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
