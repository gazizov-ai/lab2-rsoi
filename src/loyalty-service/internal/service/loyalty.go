package service

import (
	"context"

	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/model"
	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/repository"
)

type LoyaltyService struct {
	repo *repository.LoyaltyRepository
}

func NewLoyaltyService(repo *repository.LoyaltyRepository) *LoyaltyService {
	return &LoyaltyService{repo: repo}
}

func (s *LoyaltyService) Health(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *LoyaltyService) GetLoyalty(ctx context.Context, username string) (model.LoyaltyResponse, error) {
	return s.repo.GetLoyalty(ctx, username)
}

func (s *LoyaltyService) IncrementReservationCount(ctx context.Context, username string) error {
	return s.repo.IncrementReservationCount(ctx, username)
}
