package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gazizov-ai/lab2-rsoi/src/payment-service/internal/model"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment model.PaymentResponse) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO payments(payment_uid, username, status, price)
		 VALUES ($1, $2, $3, $4)`,
		payment.PaymentUID, payment.Username, payment.Status, payment.Price,
	)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}
	return nil
}

func (r *PaymentRepository) GetPayment(ctx context.Context, uid string) (model.PaymentResponse, error) {
	var p model.PaymentResponse

	err := r.db.QueryRowContext(ctx,
		`SELECT payment_uid, username, status, price
		 FROM payments WHERE payment_uid = $1`,
		uid,
	).Scan(&p.PaymentUID, &p.Username, &p.Status, &p.Price)

	if err == sql.ErrNoRows {
		return model.PaymentResponse{}, nil
	}
	if err != nil {
		return model.PaymentResponse{}, fmt.Errorf("select payment: %w", err)
	}

	return p, nil
}

func (r *PaymentRepository) CancelPayment(ctx context.Context, uid string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE payments SET status = 'CANCELED'
		 WHERE payment_uid = $1`,
		uid,
	)
	return err
}

func (r *PaymentRepository) GetPaymentsByUser(ctx context.Context, username string) ([]model.PaymentResponse, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT payment_uid, username, status, price
		 FROM payments WHERE username = $1
		 ORDER BY id`,
		username,
	)
	if err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	defer rows.Close()

	var result []model.PaymentResponse

	for rows.Next() {
		var p model.PaymentResponse
		if err := rows.Scan(
			&p.PaymentUID,
			&p.Username,
			&p.Status,
			&p.Price,
		); err != nil {
			return nil, fmt.Errorf("scan payment: %w", err)
		}
		result = append(result, p)
	}

	return result, rows.Err()
}
