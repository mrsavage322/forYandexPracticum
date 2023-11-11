package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"net/http"
	"time"
)

func BDConnection(w http.ResponseWriter, r *http.Request) {
	Database := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `videos`, `userpassword`, `videos`)

	db, err := sql.Open("pgx", Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		http.Error(w, "Non-existent identifier", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	//err = db.Ping()
	//if err != nil {
	//	log.Fatal("Ошибка при попытке установить соединение:", err)
	//}
	//fmt.Println("Соединение установлено успешно!")
}
