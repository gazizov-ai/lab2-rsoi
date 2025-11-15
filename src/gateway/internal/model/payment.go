package model

type Payment struct {
	PaymentUID string `json:"paymentUid"`
	Username   string `json:"username"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}
