package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/model"
	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/service"
)

type Handler struct {
	svc *service.ReservationService
}

func NewHandler(s *service.ReservationService) *Handler {
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

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Username   string `json:"username"`
		HotelUID   string `json:"hotelUid"`
		StartDate  string `json:"startDate"`
		EndDate    string `json:"endDate"`
		PaymentUID string `json:"paymentUid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	end, err := time.Parse("2006-01-02", body.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := model.CreateReservationRequest{
		Username:   body.Username,
		HotelUID:   body.HotelUID,
		StartDate:  start,
		EndDate:    end,
		PaymentUID: body.PaymentUID,
	}

	res, err := h.svc.CreateReservation(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	uid := last(r.URL.Path)
	res, err := h.svc.GetReservation(r.Context(), uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.ReservationUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetReservationsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	username := last(r.URL.Path)
	res, err := h.svc.GetReservationsByUser(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	uid := last(r.URL.Path)
	if err := h.svc.CancelReservation(r.Context(), uid); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListHotels(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetHotel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hotelUID := last(r.URL.Path)
	hh, err := h.svc.GetHotel(r.Context(), hotelUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if hh.HotelUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(hh)
}

func last(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
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
