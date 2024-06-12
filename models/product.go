package models

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" form:"name" binding:"required"`
	Description string  `json:"description" form:"description" binding:"required"`
	ImgUrl      string  `json:"img_url"`
	Price       float64 `json:"price" form:"price" binding:"required"`
	CategoryId  int64   `json:"category_id"`
}
