package main

type Contract struct {
	ContractId              string `json:"contractId"`
	Order                   Order  `json:"order"`
	CreateDate              string `json:"createDate"`
	ContractDueDate         string `json:"contractDueDate"`
	ContractStatus          string `json:"contractStatus"`
	AdditionalInf           string `json:"additionalInf"`
	Invoice					string `json:"invoice"`
	PackagingList			string `json:"packagingList"`
	LineOfCredit			string `json:"lineOfCredit"`
	BillOfShipment			string `json:"billOfShipment"`
}

type Order struct {
	Articles []Articles `json:"articles"`
	Amount   Amount     `json:"amount"`
	Buyer    User       `json:"buyer"`
	Seller   User       `json:"seller"`
	Shipment Shipment   `json:"shipment"`
}

type User struct {
	UniqueId string `json:"uniqueId"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Bank     string `json:"bank"`
	Address  string `json:"address"`
}

type Amount struct {
	Currency string `json:"currency"`
	Value    int    `json:"value"`
}

type Document struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type Articles struct {
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Amount      Amount `json:"amount"`
}

type Shipment struct {
	Company    Company `json:"company"`
	TrackingID string  `json:"trackingId"`
}

type Company struct {
	UniqueId string `json:"uniqueId"`
	Name     string `json:"name"`
	Country  string `json:"country"`
}
