package app

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"net/http"
)

func BDConnection(w http.ResponseWriter, r *http.Request) {
	ps := fmt.Sprintf("host=%s user=%s password=%s sslmode=mode",
		`localhost`, `video`, `XXXXXXXXX`, `video`)

	db, err := sql.Open("pgx", ps)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
