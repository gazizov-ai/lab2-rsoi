package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gazizov-ai/lab2-rsoi/src/reservation-service/internal/model"
)

type ReservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *ReservationRepository) CreateReservation(ctx context.Context, res model.Reservation) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO reservations (reservation_uid, username, hotel_id, start_date, end_data, status, payment_uid)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		res.ReservationUID,
		res.Username,
		res.HotelID,
		res.StartDate,
		res.EndDate,
		res.Status,
		res.PaymentUID,
	)
	if err != nil {
		return fmt.Errorf("insert reservation: %w", err)
	}
	return nil
}

func (r *ReservationRepository) GetReservation(ctx context.Context, uid string) (model.Reservation, error) {
	var res model.Reservation

	err := r.db.QueryRowContext(ctx, `
		SELECT r.reservation_uid, r.username, h.hotel_uid, r.hotel_id, r.start_date, r.end_data, r.status, r.payment_uid
		FROM reservations r
		JOIN hotels h ON h.id = r.hotel_id
		WHERE reservation_uid = $1
	`, uid).Scan(
		&res.ReservationUID,
		&res.Username,
		&res.HotelUID,
		&res.HotelID,
		&res.StartDate,
		&res.EndDate,
		&res.Status,
		&res.PaymentUID,
	)

	if err == sql.ErrNoRows {
		return model.Reservation{}, nil
	}
	if err != nil {
		return model.Reservation{}, fmt.Errorf("select reservation: %w", err)
	}

	return res, nil
}

func (r *ReservationRepository) GetReservationsByUser(ctx context.Context, username string) ([]model.Reservation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT r.reservation_uid, r.username, h.hotel_uid, r.hotel_id, r.start_date, r.end_data, r.status, r.payment_uid
		FROM reservations r
		JOIN hotels h ON h.id = r.hotel_id
		WHERE username = $1
		ORDER BY start_date DESC
	`, username)
	if err != nil {
		return nil, fmt.Errorf("list reservations: %w", err)
	}
	defer rows.Close()

	var res []model.Reservation
	for rows.Next() {
		var rsv model.Reservation
		if err := rows.Scan(
			&rsv.ReservationUID,
			&rsv.Username,
			&rsv.HotelUID,
			&rsv.HotelID,
			&rsv.StartDate,
			&rsv.EndDate,
			&rsv.Status,
			&rsv.PaymentUID,
		); err != nil {
			return nil, fmt.Errorf("scan reservation: %w", err)
		}
		res = append(res, rsv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return res, nil
}

func (r *ReservationRepository) CancelReservation(ctx context.Context, uid string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE reservations
		SET status = 'CANCELED'
		WHERE reservation_uid = $1
	`, uid)
	if err != nil {
		return fmt.Errorf("cancel reservation: %w", err)
	}
	return nil
}

func (r *ReservationRepository) GetHotelIDByUID(ctx context.Context, uid string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM hotels WHERE hotel_uid = $1`,
		uid,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("select hotel id: %w", err)
	}
	return id, nil
}

func (r *ReservationRepository) ListHotels(ctx context.Context, page, size int) ([]model.Hotel, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM hotels`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count hotels: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT hotel_uid, name, country, city, address, stars, price
		FROM hotels
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, size, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("select hotels: %w", err)
	}
	defer rows.Close()

	var items []model.Hotel
	for rows.Next() {
		var h model.Hotel
		if err := rows.Scan(
			&h.HotelUID,
			&h.Name,
			&h.Country,
			&h.City,
			&h.Address,
			&h.Stars,
			&h.Price,
		); err != nil {
			return nil, 0, fmt.Errorf("scan hotel: %w", err)
		}
		items = append(items, h)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return items, total, nil
}

func (r *ReservationRepository) GetHotelByUID(ctx context.Context, hotelUID string) (model.Hotel, error) {
	var h model.Hotel

	err := r.db.QueryRowContext(ctx, `
		SELECT hotel_uid, name, country, city, address, stars, price
		FROM hotels
		WHERE hotel_uid = $1
	`, hotelUID).Scan(
		&h.HotelUID,
		&h.Name,
		&h.Country,
		&h.City,
		&h.Address,
		&h.Stars,
		&h.Price,
	)

	if err == sql.ErrNoRows {
		return model.Hotel{}, nil
	}
	if err != nil {
		return model.Hotel{}, fmt.Errorf("get hotel: %w", err)
	}

	return h, nil
}
