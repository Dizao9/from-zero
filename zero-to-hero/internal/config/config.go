package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DSN  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load config file")
		return nil, err
	}

	port := os.Getenv("SERVER_PORT")
	dsn := os.Getenv("DB_DSN")

	if port == "" {
		log.Fatalf("please, set the SERVER_PORT env")
	}

	if dsn == "" {
		log.Fatalf("please, set the DB_DSN env")
	}
	return &Config{
		Port: port,
		DSN:  dsn,
	}, nil
}
