package domain

type Payment struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
}
