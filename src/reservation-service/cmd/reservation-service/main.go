package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/config"
	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/httpserver"
	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/repository"
	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping error: %v", err)
	}

	repo := repository.NewReservationRepository(db)
	svc := service.NewReservationService(repo)
	router := httpserver.NewRouter(svc)

	log.Printf("reservation-service listening on %s", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
