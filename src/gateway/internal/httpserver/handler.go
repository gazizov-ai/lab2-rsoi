package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/service"
)

type Handler struct {
	svc service.Gateway
}

func NewHandler(s service.Gateway) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if err := h.svc.Health(r.Context()); err != nil {
		WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unhealthy"})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func getUsername(r *http.Request) string {
	return r.Header.Get("X-User-Name")
}

func (h *Handler) Hotels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := r.URL.Query()
	page := parseIntOrDefault(q.Get("page"), 1)
	size := parseIntOrDefault(q.Get("size"), 10)

	resp, err := h.svc.ListHotels(r.Context(), page, size)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) Loyalty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	username := getUsername(r)
	if username == "" {
		WriteError(w, http.StatusUnauthorized, "missing X-User-Name header")
		return
	}

	resp, err := h.svc.GetLoyalty(username)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) ListReservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	username := getUsername(r)
	if username == "" {
		WriteError(w, http.StatusUnauthorized, "missing X-User-Name header")
		return
	}
	resp, err := h.svc.ListUserReservations(r.Context(), username)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	username := getUsername(r)
	if username == "" {
		WriteError(w, http.StatusUnauthorized, "missing X-User-Name header")
		return
	}

	var body struct {
		HotelUID  string `json:"hotelUid"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	resp, err := h.svc.CreateReservation(r.Context(), username, body.HotelUID, body.StartDate, body.EndDate)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	username := getUsername(r)
	if username == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reservationUID := last(r.URL.Path)
	if reservationUID == "" {
		WriteError(w, http.StatusBadRequest, "invalid reservation uid")
		return
	}

	resp, err := h.svc.GetReservation(r.Context(), username, reservationUID)

	if err != nil {
		if err.Error() == "forbidden" {
			WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.ReservationUID == "" {
		WriteError(w, http.StatusNotFound, "not found")
		return
	}
	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	username := getUsername(r)
	if username == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reservationUID := last(r.URL.Path)
	if reservationUID == "" {
		WriteError(w, http.StatusBadRequest, "invalid reservation uid")
		return
	}

	if err := h.svc.CancelReservation(r.Context(), username, reservationUID); err != nil {
		if err.Error() == "forbidden" {
			WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	username := getUsername(r)
	if username == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.svc.Me(r.Context(), username)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, resp)
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
