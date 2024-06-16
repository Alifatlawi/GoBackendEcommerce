package models

type Order struct {
	ID          int    `json:"id"`
	ProductID   int    `json:"product_id"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}
