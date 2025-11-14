package httpserver

import (
	"net/http"

	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/service"
)

func NewRouter(s *service.LoyaltyService) http.Handler {
	mux := http.NewServeMux()

	h := NewHandler(s)

	mux.HandleFunc("/manage/health", h.Health)
	mux.HandleFunc("/internal/loyalty/", h.Loyalty)

	return mux
}
