package config

import (
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string

	ReservationURL string
	PaymentURL     string
	LoyaltyURL     string
}

func Load() Config {
	return Config{
		Host:           getenv("HOST", "0.0.0.0"),
		Port:           getenv("PORT", "8080"),
		ReservationURL: getenv("RESERVATION_URL", "http://reservation-service:8070"),
		PaymentURL:     getenv("PAYMENT_URL", "http://payment-service:8060"),
		LoyaltyURL:     getenv("LOYALTY_URL", "http://loyalty-service:8050"),
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
