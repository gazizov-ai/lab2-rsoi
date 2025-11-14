package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/service"
)

type Handler struct {
	svc *service.GatewayService
}

func NewHandler(s *service.GatewayService) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := h.svc.Health(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func getUsername(r *http.Request) string {
	return r.Header.Get("X-User-Name")
}

func (h *Handler) Hotels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	page := parseIntOrDefault(q.Get("page"), 1)
	size := parseIntOrDefault(q.Get("size"), 10)

	resp, err := h.svc.ListHotels(r.Context(), page, size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Loyalty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := getUsername(r)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := h.svc.GetLoyalty(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListReservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := getUsername(r)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	resp, err := h.svc.ListUserReservations(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := getUsername(r)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var body struct {
		HotelUID  string `json:"hotelUid"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.svc.CreateReservation(r.Context(), username, body.HotelUID, body.StartDate, body.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := getUsername(r)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reservationUID := last(r.URL.Path)
	resp, err := h.svc.GetReservation(r.Context(), username, reservationUID)
	if err != nil {
		if err.Error() == "forbidden" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.ReservationUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := getUsername(r)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reservationUID := last(r.URL.Path)
	if err := h.svc.CancelReservation(r.Context(), username, reservationUID); err != nil {
		if err.Error() == "forbidden" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := getUsername(r)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := h.svc.Me(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func parseIntOrDefault(raw string, def int) int {
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

func last(path string) string {
	n := len(path)
	if n == 0 {
		return ""
	}
	i := n - 1
	for i >= 0 && path[i] != '/' {
		i--
	}
	return path[i+1:]
}
