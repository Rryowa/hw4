package models

import "time"

type Config struct {
	User     string        `env:"POSTGRES_USER"`
	Password string        `env:"POSTGRES_PASSWORD"`
	Host     string        `env:"DB_HOST"`
	Port     string        `env:"DB_PORT"`
	DBName   string        `env:"POSTGRES_DB"`
	Attempts int           `env:"ATTEMPTS"`
	Timeout  time.Duration `env:"TIMEOUT"`
}
