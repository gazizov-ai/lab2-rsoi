package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/service"
)

type Handler struct {
	loyaltyService *service.LoyaltyService
}

func NewHandler(s *service.LoyaltyService) *Handler {
	return &Handler{loyaltyService: s}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := h.loyaltyService.Health(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) Loyalty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := parts[len(parts)-1]

	resp, err := h.loyaltyService.GetLoyalty(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "internal error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
