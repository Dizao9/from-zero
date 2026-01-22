package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Handler struct {
	DB *sql.DB
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if method := r.Method; method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rows, err := h.DB.Query(`SELECT username, email
	FROM users`)
	if err != nil {
		http.Error(w, "Internal Server Error, Failed to get users data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	res := make([]User, 0, 3)

	for rows.Next() {
		var u User
		err := rows.Scan(&u.Username, &u.Email)
		if err != nil {
			http.Error(w, "failed to read data from rows", http.StatusInternalServerError)
			return
		}
		res = append(res, u)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "error with rows method", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("failed to encode response")
	}
}

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

	h := Handler{DB: db}

	http.HandleFunc("/users", h.GetUsers)

	fmt.Println("server is running on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Println("error to start server on port 8082:", err)
	}
}
