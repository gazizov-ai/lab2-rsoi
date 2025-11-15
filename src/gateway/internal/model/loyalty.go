package model

type Loyalty struct {
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
	ReservationCount int    `json:"reservationCount"`
}
