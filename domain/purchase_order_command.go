package domain

type PurchaseOrderCheckedOut struct {
	Id         string
	ModuleName string
	ExternalId string
	CustomerId string
	ProductId  string
	Address    Address
	Payment    Payment
}

type Address struct {
	Cep    string
	Number string
}

type Payment struct {
	Token string
	Price Price
}

type Price struct {
	Amount   int
	Currency string
	Scale    int
}
