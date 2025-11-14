package service

import (
	"context"
	"errors"
	"time"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/clients"
	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/model"
)

type GatewayService struct {
	reservationClient *clients.ReservationClient
	paymentClient     *clients.PaymentClient
	loyaltyClient     *clients.LoyaltyClient
}

func NewGatewayService(
	resClient *clients.ReservationClient,
	payClient *clients.PaymentClient,
	loyalClient *clients.LoyaltyClient,
) *GatewayService {
	return &GatewayService{
		reservationClient: resClient,
		paymentClient:     payClient,
		loyaltyClient:     loyalClient,
	}
}

func (s *GatewayService) Health(ctx context.Context) error {
	return nil
}

func (s *GatewayService) ListHotels(ctx context.Context, page, size int) (model.HotelsPage, error) {
	return s.reservationClient.ListHotels(page, size)
}

func (s *GatewayService) GetLoyalty(username string) (model.Loyalty, error) {
	return s.loyaltyClient.GetLoyalty(username)
}

func (s *GatewayService) ListUserReservations(ctx context.Context, username string) ([]model.ReservationShort, error) {
	reservations, err := s.reservationClient.GetReservationsByUser(username)
	if err != nil {
		return nil, err
	}

	var result []model.ReservationShort
	for _, r := range reservations {
		h, err := s.reservationClient.GetHotel(r.HotelUID)
		if err != nil {
			return nil, err
		}
		p, err := s.paymentClient.GetPayment(r.PaymentUID)
		if err != nil {
			return nil, err
		}

		result = append(result, model.ReservationShort{
			ReservationUID: r.ReservationUID,
			Hotel:          h,
			StartDate:      r.StartDate,
			EndDate:        r.EndDate,
			Status:         r.Status,
			Payment:        p,
		})
	}

	return result, nil
}

func (s *GatewayService) GetReservation(ctx context.Context, username, reservationUID string) (model.ReservationShort, error) {
	r, err := s.reservationClient.GetReservation(reservationUID)
	if err != nil {
		return model.ReservationShort{}, err
	}
	if r.ReservationUID == "" {
		return model.ReservationShort{}, nil
	}
	if r.Username != username {
		return model.ReservationShort{}, errors.New("forbidden")
	}

	h, err := s.reservationClient.GetHotel(r.HotelUID)
	if err != nil {
		return model.ReservationShort{}, err
	}
	p, err := s.paymentClient.GetPayment(r.PaymentUID)
	if err != nil {
		return model.ReservationShort{}, err
	}

	return model.ReservationShort{
		ReservationUID: r.ReservationUID,
		Hotel:          h,
		StartDate:      r.StartDate,
		EndDate:        r.EndDate,
		Status:         r.Status,
		Payment:        p,
	}, nil
}

func (s *GatewayService) CreateReservation(ctx context.Context, username, hotelUID, startDateStr, endDateStr string) (model.ReservationShort, error) {
	hotel, err := s.reservationClient.GetHotel(hotelUID)
	if err != nil {
		return model.ReservationShort{}, err
	}
	if hotel.HotelUID == "" {
		return model.ReservationShort{}, errors.New("hotel not found")
	}

	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return model.ReservationShort{}, err
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return model.ReservationShort{}, err
	}

	days := int(end.Sub(start).Hours() / 24)
	if days <= 0 {
		days = 1
	}
	totalPrice := hotel.Price * days

	payment, err := s.paymentClient.CreatePayment(username, totalPrice)
	if err != nil {
		return model.ReservationShort{}, err
	}

	internalReq := model.ReservationInternal{
		Username:   username,
		HotelUID:   hotel.HotelUID,
		StartDate:  start,
		EndDate:    end,
		Status:     "PAID",
		PaymentUID: payment.PaymentUID,
	}

	reservation, err := s.reservationClient.CreateReservation(internalReq)
	if err != nil {
		return model.ReservationShort{}, err
	}

	return model.ReservationShort{
		ReservationUID: reservation.ReservationUID,
		Hotel:          hotel,
		StartDate:      reservation.StartDate,
		EndDate:        reservation.EndDate,
		Status:         reservation.Status,
		Payment:        payment,
	}, nil
}

func (s *GatewayService) CancelReservation(ctx context.Context, username, reservationUID string) error {
	r, err := s.reservationClient.GetReservation(reservationUID)
	if err != nil {
		return err
	}
	if r.ReservationUID == "" {
		return nil
	}
	if r.Username != username {
		return errors.New("forbidden")
	}

	if err := s.reservationClient.CancelReservation(reservationUID); err != nil {
		return err
	}
	if err := s.paymentClient.CancelPayment(r.PaymentUID); err != nil {
		return err
	}

	return nil
}

func (s *GatewayService) Me(ctx context.Context, username string) (model.MeResponse, error) {
	loyalty, err := s.loyaltyClient.GetLoyalty(username)
	if err != nil {
		return model.MeResponse{}, err
	}
	reservations, err := s.ListUserReservations(ctx, username)
	if err != nil {
		return model.MeResponse{}, err
	}

	return model.MeResponse{
		Username:     username,
		Loyalty:      loyalty,
		Reservations: reservations,
	}, nil
}
