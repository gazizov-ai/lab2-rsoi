package model

import "time"

type ReservationShort struct {
	ReservationUID string  `json:"reservationUid"`
	Hotel          Hotel   `json:"hotel"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	Status         string  `json:"status"`
	Payment        Payment `json:"payment"`
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

type ReservationFull struct {
	ReservationUID string    `json:"reservationUid"`
	Username       string    `json:"username"`
	HotelUID       string    `json:"hotelUid"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Status         string    `json:"status"`
	PaymentUID     string    `json:"paymentUid"`
	Hotel          Hotel     `json:"hotel"`
	Payment        Payment   `json:"payment"`
}

type ReservationCreateResponse struct {
	ReservationUID string                `json:"reservationUid"`
	HotelUID       string                `json:"hotelUid"`
	StartDate      string                `json:"startDate"`
	EndDate        string                `json:"endDate"`
	Discount       int                   `json:"discount"`
	Status         string                `json:"status"`
	Payment        PaymentCreateResponse `json:"payment"`
}

type PaymentCreateResponse struct {
	Status string `json:"status"`
	Price  int    `json:"price"`
}
