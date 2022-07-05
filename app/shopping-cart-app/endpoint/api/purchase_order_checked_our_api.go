package api

type PurchaseOrderCheckedOutRequest struct {
	CustomerId string `json:"customerId"`
	ProductId  string `json:"productId"`
	Address    struct {
		Cep    string `json:"cep"`
		Number string `json:"number"`
	} `json:"address"`
	Payment struct {
		Token string `json:"token"`
		Price struct {
			Amount   int    `json:"amount"`
			Currency string `json:"currency"`
			Scale    int    `json:"scale"`
		} `json:"price"`
	} `json:"payment"`
}

type PurchaseOrderCheckedOutResponse struct {
	PurchaseOrderId string `json:"id"`
}
