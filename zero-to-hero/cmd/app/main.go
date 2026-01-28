package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zero-to-hero/internal/config"
	"zero-to-hero/internal/storage"
	"zero-to-hero/internal/transport"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func connectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed open driver for connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed connect to db: %v", err)
	}
	return db, err
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config data: %v", err)
	}

	db, err := connectToDB(conf.DSN)
	if err != nil {
		log.Fatalf("failed connect to db: %v", err)
	}

	store := &storage.Storage{DB: db}
	h := &transport.Handler{Store: store}

	finalHandler := http.HandlerFunc(h.GetUsers)

	wrappedHandler := transport.LoggindMiddleware(finalHandler)
	http.Handle("/users", wrappedHandler)

	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: nil,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	db.Close()
	log.Println("Server exited properlsy")
}
