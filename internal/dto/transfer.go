package dto

type Transfer struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
	Username string `json:"username"`
}
