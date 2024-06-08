package main

type Receipt struct {
	TransactionID string `json:"transaction_id"`
	ExternalID    string `json:"external_id"`
	Total         int64  `json:"total"`
	Currency      string `json:"currency"`
	Items         []Item `json:"items"`
	Taxes         []Tax  `json:"taxes"`
}

type Item struct {
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	SubItems    []SubItem `json:"sub_items"`
}

type SubItem struct {
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Taxes       int    `json:"tax"`
}

type Tax struct {
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	TaxNumber   string `json:"tax_number"`
}
