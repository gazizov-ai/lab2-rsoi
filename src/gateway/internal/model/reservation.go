package model

import "time"

type ReservationShort struct {
	ReservationUID string    `json:"reservationUid"`
	Hotel          Hotel     `json:"hotel"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Status         string    `json:"status"`
	Payment        Payment   `json:"payment"`
}

type ReservationInternal struct {
	ReservationUID string    `json:"reservationUid"`
	Username       string    `json:"username"`
	HotelUID       string    `json:"hotelUid"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Status         string    `json:"status"`
	PaymentUID     string    `json:"paymentUid"`
}
