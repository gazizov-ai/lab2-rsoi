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
		`SELECT status, discount FROM loyalties WHERE username = $1`,
		username,
	).Scan(&resp.Status, &resp.Discount)

	if err == sql.ErrNoRows {
		return model.LoyaltyResponse{
			Status:   "Bronze",
			Discount: 5,
		}, nil
	}

	if err != nil {
		return model.LoyaltyResponse{}, fmt.Errorf("query loyalty: %w", err)
	}

	return resp, nil
}
