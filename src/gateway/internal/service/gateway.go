package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/clients"
	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/model"
)

type Gateway interface {
	Health(ctx context.Context) error
	ListHotels(ctx context.Context, page, size int) (model.HotelsPage, error)
	GetLoyalty(username string) (model.Loyalty, error)
	ListUserReservations(ctx context.Context, username string) ([]model.ReservationShort, error)
	GetReservation(ctx context.Context, username, reservationUID string) (model.ReservationShort, error)
	CreateReservation(ctx context.Context, username, hotelUID, startDateStr, endDateStr string) (model.ReservationCreateResponse, error)
	CancelReservation(ctx context.Context, username, reservationUID string) error
	Me(ctx context.Context, username string) (model.MeResponse, error)
}

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

		h.FullAddress = fmt.Sprintf("%s, %s, %s", h.Country, h.City, h.Address)

		result = append(result, model.ReservationShort{
			ReservationUID: r.ReservationUID,
			Hotel:          h,
			StartDate:      r.StartDate.Format("2006-01-02"),
			EndDate:        r.EndDate.Format("2006-01-02"),
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
	h.FullAddress = fmt.Sprintf("%s, %s, %s", h.Country, h.City, h.Address)
	if err != nil {
		return model.ReservationShort{}, err
	}

	return model.ReservationShort{
		ReservationUID: r.ReservationUID,
		Hotel:          h,
		StartDate:      r.StartDate.Format("2006-01-02"),
		EndDate:        r.EndDate.Format("2006-01-02"),
		Status:         r.Status,
		Payment:        p,
	}, nil
}

func (s *GatewayService) CreateReservation(ctx context.Context, username, hotelUID, startDateStr, endDateStr string) (model.ReservationCreateResponse, error) {
	hotel, err := s.reservationClient.GetHotel(hotelUID)
	if err != nil {
		return model.ReservationCreateResponse{}, err
	}
	if hotel.HotelUID == "" {
		return model.ReservationCreateResponse{}, errors.New("hotel not found")
	}

	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return model.ReservationCreateResponse{}, err
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return model.ReservationCreateResponse{}, err
	}

	loyalty, err := s.loyaltyClient.GetLoyalty(username)
	if err != nil {
		return model.ReservationCreateResponse{}, err
	}

	days := int(end.Sub(start).Hours() / 24)
	if days <= 0 {
		days = 1
	}
	basePrice := hotel.Price * days
	finalPrice := basePrice - (basePrice * loyalty.Discount / 100)

	payment, err := s.paymentClient.CreatePayment(username, finalPrice)
	if err != nil {
		return model.ReservationCreateResponse{}, err
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
		return model.ReservationCreateResponse{}, err
	}

	if err := s.loyaltyClient.IncrementReservation(username); err != nil {
		return model.ReservationCreateResponse{}, err
	}

	resp := model.ReservationCreateResponse{
		ReservationUID: reservation.ReservationUID,
		HotelUID:       hotel.HotelUID,
		StartDate:      startDateStr,
		EndDate:        endDateStr,
		Discount:       loyalty.Discount,
		Status:         reservation.Status,
		Payment: model.PaymentCreateResponse{
			Status: payment.Status,
			Price:  finalPrice,
		},
	}

	return resp, nil
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
