package httpserver

import (
	"net/http"

	"github.com/gazizov-ai/lab2-rsoi/src/payment-service/internal/service"
)

func NewRouter(s *service.PaymentService) http.Handler {
	mux := http.NewServeMux()

	h := NewHandler(s)

	mux.HandleFunc("/manage/health", h.Health)

	mux.HandleFunc("/internal/payments", h.CreatePayment)
	mux.HandleFunc("/internal/payments/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetPayment(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			h.CancelPayment(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/internal/payments/byUser/", h.GetPaymentsByUser)

	return mux
}
