package httpserver

import (
	"net/http"

	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/service"
)

func NewRouter(s *service.ReservationService) http.Handler {
	mux := http.NewServeMux()

	h := NewHandler(s)

	mux.HandleFunc("/manage/health", h.Health)

	mux.HandleFunc("/internal/reservations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateReservation(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/internal/reservations/byUser/", h.GetReservationsByUser)

	mux.HandleFunc("/internal/reservations/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetReservation(w, r)
		case http.MethodDelete:
			h.CancelReservation(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/internal/hotels", h.ListHotels)
	mux.HandleFunc("/internal/hotels/", h.GetHotel)

	return mux
}
