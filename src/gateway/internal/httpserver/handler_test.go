package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/model"
)

type fakeGateway struct {
	healthErr         error
	hotelsPage        model.HotelsPage
	loyalty           model.Loyalty
	reservations      []model.ReservationShort
	me                model.MeResponse
	getReservationRes model.ReservationShort
	getReservationErr error
}

func (f *fakeGateway) Health(_ context.Context) error {
	return f.healthErr
}

func (f *fakeGateway) ListHotels(_ context.Context, page, size int) (model.HotelsPage, error) {
	return f.hotelsPage, nil
}

func (f *fakeGateway) GetLoyalty(username string) (model.Loyalty, error) {
	return f.loyalty, nil
}

func (f *fakeGateway) ListUserReservations(_ context.Context, username string) ([]model.ReservationShort, error) {
	return f.reservations, nil
}

func (f *fakeGateway) GetReservation(_ context.Context, username, reservationUID string) (model.ReservationShort, error) {
	return f.getReservationRes, f.getReservationErr
}

func (f *fakeGateway) CreateReservation(_ context.Context, username, hotelUID, startDateStr, endDateStr string) (model.ReservationCreateResponse, error) {
	// Для этих тестов нам create не важен, можно вернуть заглушку.
	return model.ReservationCreateResponse{}, nil
}

func (f *fakeGateway) CancelReservation(_ context.Context, username, reservationUID string) error {
	return nil
}

func (f *fakeGateway) Me(_ context.Context, username string) (model.MeResponse, error) {
	return f.me, nil
}

// helper для чтения JSON-ответа в тестах
func decodeJSONBody(t *testing.T, rr *httptest.ResponseRecorder, dst interface{}) {
	t.Helper()
	if err := json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(dst); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
}

func TestHotels_OK(t *testing.T) {
	fake := &fakeGateway{
		hotelsPage: model.HotelsPage{
			Items: []model.Hotel{
				{
					HotelUID:    "049161bb-badd-4fa8-9d90-87c9a82b0668",
					Name:        "Ararat Park Hyatt Moscow",
					Country:     "Россия",
					City:        "Москва",
					Address:     "Неглинная ул., 4",
					Stars:       5,
					Price:       10000,
					FullAddress: "Россия, Москва, Неглинная ул., 4",
				},
			},
			TotalElements: 1,
		},
	}
	h := NewHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hotels?page=1&size=10", nil)
	rr := httptest.NewRecorder()

	h.Hotels(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp model.HotelsPage
	decodeJSONBody(t, rr, &resp)

	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 hotel, got %d", len(resp.Items))
	}
	if resp.Items[0].HotelUID != fake.hotelsPage.Items[0].HotelUID {
		t.Fatalf("unexpected hotelUid: %s", resp.Items[0].HotelUID)
	}
}

func TestLoyalty_UnauthorizedWithoutHeader(t *testing.T) {
	fake := &fakeGateway{}
	h := NewHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/loyalty", nil)
	rr := httptest.NewRecorder()

	h.Loyalty(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestLoyalty_OK(t *testing.T) {
	fake := &fakeGateway{
		loyalty: model.Loyalty{
			Status:           "GOLD",
			Discount:         10,
			ReservationCount: 26,
		},
	}
	h := NewHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/loyalty", nil)
	req.Header.Set("X-User-Name", "Test Max")
	rr := httptest.NewRecorder()

	h.Loyalty(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp model.Loyalty
	decodeJSONBody(t, rr, &resp)

	if resp.Status != "GOLD" || resp.Discount != 10 || resp.ReservationCount != 26 {
		t.Fatalf("unexpected loyalty response: %+v", resp)
	}
}

func TestMe_OK(t *testing.T) {
	fake := &fakeGateway{
		me: model.MeResponse{
			Username: "Test Max",
			Loyalty: model.Loyalty{
				Status:           "GOLD",
				Discount:         10,
				ReservationCount: 26,
			},
			Reservations: []model.ReservationShort{
				{
					ReservationUID: "e2866665-68f0-464b-802f-3a6eae827895",
					Hotel: model.Hotel{
						HotelUID:    "049161bb-badd-4fa8-9d90-87c9a82b0668",
						Name:        "Ararat Park Hyatt Moscow",
						Country:     "Россия",
						City:        "Москва",
						Address:     "Неглинная ул., 4",
						Stars:       5,
						Price:       10000,
						FullAddress: "Россия, Москва, Неглинная ул., 4",
					},
					StartDate: "2021-10-08",
					EndDate:   "2021-10-11",
					Status:    "PAID",
					Payment: model.Payment{
						PaymentUID: "11111111-1111-1111-1111-111111111111",
						Username:   "Test Max",
						Status:     "PAID",
						Price:      27000,
					},
				},
			},
		},
	}
	h := NewHandler(fake)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("X-User-Name", "Test Max")
	rr := httptest.NewRecorder()

	h.Me(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp model.MeResponse
	decodeJSONBody(t, rr, &resp)

	if resp.Username != "Test Max" {
		t.Fatalf("unexpected username: %s", resp.Username)
	}
	if len(resp.Reservations) != 1 {
		t.Fatalf("expected 1 reservation, got %d", len(resp.Reservations))
	}
	if resp.Reservations[0].Hotel.FullAddress != "Россия, Москва, Неглинная ул., 4" {
		t.Fatalf("unexpected fullAddress: %s", resp.Reservations[0].Hotel.FullAddress)
	}
}
