package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"zero-to-hero/internal/storage"
	"zero-to-hero/internal/transport"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func connectToDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open a db driver:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	return db
}

func main() {
	dsn := "postgres://postgres:postgres@localhost:5433/postgres"
	db := connectToDB(dsn)
	defer db.Close()

	s := &storage.Storage{DB: db}
	h := &transport.Handler{Store: s}

	http.HandleFunc("/users", h.GetUsers)

	fmt.Println("server is running on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Println("error to start server on port 8082:", err)
	}
}
