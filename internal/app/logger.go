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

		next.ServeHTTP(w, r)
		elapsed := time.Since(startTime)

		sugarLogger.Infow("Request",
			"URL", r.URL,
			"Method", r.Method,
			"Elapsed", elapsed,
		)
	})
}
