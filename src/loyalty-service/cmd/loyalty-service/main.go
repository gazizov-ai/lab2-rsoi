package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/config"
	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/httpserver"
	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/repository"
	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/service"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	repo := repository.NewLoyaltyRepository(db)
	svc := service.NewLoyaltyService(repo)
	router := httpserver.NewRouter(svc)

	log.Printf("loyalty-service listening on %s", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
