package model

type MeResponse struct {
	Username     string             `json:"username"`
	Loyalty      Loyalty            `json:"loyalty"`
	Reservations []ReservationShort `json:"reservations"`
}
