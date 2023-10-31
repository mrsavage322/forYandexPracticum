package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/mrsavage322/foryandex/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app.URLMap = app.NewURLMapStorage()
	app.SetFlags()
	app.SetConfig()
	app.InitializeLogger()

	r := chi.NewRouter()
	r.Use(app.LogRequest)
	r.Use(app.GzipMiddleware)
	r.Get("/", app.Redirect)
	r.Get("/{id}", app.Redirect)
	r.Post("/", app.HandlePost)
	r.Post("/api/shorten", app.HandleJSON)

	srv := &http.Server{
		Addr:    app.ServerAddr,
		Handler: r,
	}

	// Создаем канал для сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit

		log.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Server is listening on %s\n", app.ServerAddr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Server has stopped.")
}
