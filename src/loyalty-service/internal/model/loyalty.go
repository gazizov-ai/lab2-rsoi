package model

type LoyaltyResponse struct {
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
	ReservationCount int    `json:"reservationCount"`
}
