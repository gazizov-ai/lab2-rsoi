package httpserver

import (
	"net/http"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/service"
)

func NewRouter(s service.Gateway) http.Handler {
	mux := http.NewServeMux()

	h := NewHandler(s)

	mux.HandleFunc("/manage/health", h.Health)

	mux.HandleFunc("/api/v1/hotels", h.Hotels)
	mux.HandleFunc("/api/v1/loyalty", h.Loyalty)
	mux.HandleFunc("/api/v1/reservations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ListReservations(w, r)
			return
		}
		if r.Method == http.MethodPost {
			h.CreateReservation(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/api/v1/reservations/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetReservation(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			h.CancelReservation(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/api/v1/me", h.Me)

	return mux
}
