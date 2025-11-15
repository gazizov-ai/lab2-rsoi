package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/model"
	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/repository"
)

type ReservationService struct {
	repo *repository.ReservationRepository
}

func NewReservationService(repo *repository.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) Health(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *ReservationService) CreateReservation(ctx context.Context, req model.CreateReservationRequest) (model.Reservation, error) {
	hotelID, err := s.repo.GetHotelIDByUID(ctx, req.HotelUID)
	if err != nil {
		return model.Reservation{}, err
	}
	if hotelID == 0 {
		return model.Reservation{}, fmt.Errorf("hotel not found")
	}

	res := model.Reservation{
		ReservationUID: uuid.New().String(),
		Username:       req.Username,
		HotelUID:       req.HotelUID,
		HotelID:        hotelID,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         "PAID",
		PaymentUID:     req.PaymentUID,
	}

	if err := s.repo.CreateReservation(ctx, res); err != nil {
		return model.Reservation{}, err
	}

	return res, nil
}

func (s *ReservationService) GetReservation(ctx context.Context, uid string) (model.Reservation, error) {
	return s.repo.GetReservation(ctx, uid)
}

func (s *ReservationService) GetReservationsByUser(ctx context.Context, username string) ([]model.Reservation, error) {
	return s.repo.GetReservationsByUser(ctx, username)
}

func (s *ReservationService) CancelReservation(ctx context.Context, uid string) error {
	return s.repo.CancelReservation(ctx, uid)
}

func (s *ReservationService) ListHotels(ctx context.Context, page, size int) (model.HotelsPage, error) {
	items, total, err := s.repo.ListHotels(ctx, page, size)
	if err != nil {
		return model.HotelsPage{}, err
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	return model.HotelsPage{
		Page:          page,
		PageSize:      size,
		TotalElements: total,
		Items:         items,
	}, nil
}

func (s *ReservationService) GetHotel(ctx context.Context, hotelUID string) (model.Hotel, error) {
	return s.repo.GetHotelByUID(ctx, hotelUID)
}
