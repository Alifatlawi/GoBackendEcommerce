package models

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImgUrl      string  `json:"img_url"`
	CategoryId  int     `json:"category_id"`
	Price       float64 `json:"price"`
}
