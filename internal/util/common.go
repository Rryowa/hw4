package util

import (
	"github.com/joho/godotenv"
	"homework/internal/models"
	"log"
	"os"
	"strconv"
	"time"
)

func NewConfig() *models.Config {
	err := godotenv.Load("internal/util/.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	attempts, err := strconv.Atoi(os.Getenv("ATTEMPTS"))
	if err != nil {
		log.Fatalf("err converting ATTEMPTS: %v\n", err)
	}

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v\n", err)
	}

	return &models.Config{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Attempts: attempts,
		Timeout:  timeout,
	}
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}
