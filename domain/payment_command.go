package domain

type PaymentCmd struct {
	Id             string `json:"id"`
	PaymentAttempt `json:"paymentAttempt"`
}

type PaymentAttempt struct {
	PaymentId string `json:"paymentId"`
	Status    string `json:"status"`
	Price     Price  `json:"price"`
}
