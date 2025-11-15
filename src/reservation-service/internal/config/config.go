package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port  string
	DBDSN string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8070"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://program:test@localhost:5432/loyalty?sslmode=disable"
	}

	return Config{
		Port:  port,
		DBDSN: dsn,
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%s", c.Port)
}
