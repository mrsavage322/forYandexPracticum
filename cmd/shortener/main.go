package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/mrsavage322/foryandex/internal/app"
	"github.com/mrsavage322/foryandex/internal/app/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var once sync.Once

func main() {
	app.SetFlags()
	app.SetConfig()
	app.Cfg.URLMap = app.NewURLMapStorage()
	once.Do(func() {
		app.Cfg.URLMapDB = app.NewURLDBStorage(app.Cfg.DatabaseAddr)
	})
	app.InitializeLogger()

	r := chi.NewRouter()
	r.Use(app.LogRequest)
	r.Use(handler.GzipMiddleware)
	r.Use(app.AuthMiddleware)
	r.Get("/", handler.Redirect)
	r.Get("/{id}", handler.Redirect)
	r.Get("/ping", handler.BDConnection)
	r.Get("/api/user/urls", handler.GetUserURLs)
	r.Post("/", handler.HandlePost)
	r.Post("/api/shorten", handler.HandleJSON)
	r.Post("/api/shorten/batch", handler.HandleBatch)
	r.Delete("/api/user/urls", handler.DeleteURLsHandler)

	srv := &http.Server{
		Addr:    app.Cfg.ServerAddr,
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

	log.Printf("Server is listening on %s\n", app.Cfg.ServerAddr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Server has stopped.")
}

//comment
