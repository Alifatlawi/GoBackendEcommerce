package models

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	ImgUrl      string `json:"img_url"`
	Price       string `json:"price" form:"price" binding:"required"`
	CategoryId  string `json:"category_id"`
}
