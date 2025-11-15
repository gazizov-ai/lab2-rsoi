package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/gazizov-ai/lab2-rsoi/src/payment-service/internal/model"
	"github.com/gazizov-ai/lab2-rsoi/src/payment-service/internal/repository"
)

type PaymentService struct {
	repo *repository.PaymentRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Health(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *PaymentService) CreatePayment(ctx context.Context, username string, price int) (model.PaymentResponse, error) {
	p := model.PaymentResponse{
		PaymentUID: uuid.New().String(),
		Username:   username,
		Status:     "PAID",
		Price:      price,
	}

	if err := s.repo.CreatePayment(ctx, p); err != nil {
		return model.PaymentResponse{}, err
	}

	return p, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, uid string) (model.PaymentResponse, error) {
	return s.repo.GetPayment(ctx, uid)
}

func (s *PaymentService) CancelPayment(ctx context.Context, uid string) error {
	return s.repo.CancelPayment(ctx, uid)
}

func (s *PaymentService) GetPaymentsByUser(ctx context.Context, username string) ([]model.PaymentResponse, error) {
	return s.repo.GetPaymentsByUser(ctx, username)
}
