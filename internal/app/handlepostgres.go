package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

func BDConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(), DatabaseAddr)
	if err != nil {
		fmt.Println("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = conn.Ping(ctx); err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
