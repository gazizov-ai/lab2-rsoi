package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gazizov-ai/lab2-rsoi/src/loyalty-service/internal/model"
)

type LoyaltyRepository struct {
	db *sql.DB
}

func NewLoyaltyRepository(db *sql.DB) *LoyaltyRepository {
	return &LoyaltyRepository{db: db}
}

func (r *LoyaltyRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *LoyaltyRepository) GetLoyalty(ctx context.Context, username string) (model.LoyaltyResponse, error) {
	var resp model.LoyaltyResponse

	err := r.db.QueryRowContext(ctx,
		`SELECT status, discount, reservation_count FROM loyalties WHERE username = $1`,
		username,
	).Scan(&resp.Status, &resp.Discount, &resp.ReservationCount)

	if err == sql.ErrNoRows {
		return model.LoyaltyResponse{
			Status:           "Bronze",
			Discount:         5,
			ReservationCount: 0,
		}, nil
	}

	if err != nil {
		return model.LoyaltyResponse{}, fmt.Errorf("query loyalty: %w", err)
	}

	return resp, nil
}

func (r *LoyaltyRepository) IncrementReservationCount(ctx context.Context, username string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE loyalties SET reservation_count = reservation_count + 1 WHERE username = $1`,
		username,
	)
	if err != nil {
		return fmt.Errorf("increment reservation_count: %w", err)
	}
	return nil
}
