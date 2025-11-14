package main

import (
	"log"
	"net/http"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/clients"
	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/config"
	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/httpserver"
	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/service"
)

func main() {
	cfg := config.Load()

	resClient := clients.NewReservationClient(cfg.ReservationURL)
	payClient := clients.NewPaymentClient(cfg.PaymentURL)
	loyalClient := clients.NewLoyaltyClient(cfg.LoyaltyURL)

	svc := service.NewGatewayService(resClient, payClient, loyalClient)
	router := httpserver.NewRouter(svc)

	log.Printf("gateway listening on %s", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
