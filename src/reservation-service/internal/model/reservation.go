package model

import "time"

type Reservation struct {
	ReservationUID string    `json:"reservationUid"`
	Username       string    `json:"username"`
	HotelUID       string    `json:"hotelUid"`
	HotelID        int       `json:"-"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Status         string    `json:"status"`
	PaymentUID     string    `json:"paymentUid"`
}

type CreateReservationRequest struct {
	Username   string    `json:"username"`
	HotelUID   string    `json:"hotelUid"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	PaymentUID string    `json:"paymentUid"`
}
