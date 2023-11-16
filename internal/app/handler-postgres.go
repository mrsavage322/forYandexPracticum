package app

import (
	"context"
	"github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

func BDConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(), Cfg.DatabaseAddr)
	if err != nil {
		sugar.Error("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = conn.Ping(ctx); err != nil {
		sugar.Error("Failed to connect to the database:", err)
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
