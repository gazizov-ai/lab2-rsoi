package model

type PaymentResponse struct {
	PaymentUID string `json:"paymentUid"`
	Username   string `json:"username"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}
