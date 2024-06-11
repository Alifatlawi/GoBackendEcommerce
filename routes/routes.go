package routes

import (
	"ecommercebackend/controllers"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/categories", controllers.GetCategories)
		api.POST("/categories", controllers.CreateCategory)
		api.GET("/products", controllers.GetProducts)
		api.POST("/products", controllers.CreateProduct)
		api.PUT("/products", controllers.UpdateProduct)
		api.DELETE("/products", controllers.DeleteProduct)
		api.GET("/product/:id", controllers.GetProductById)
		api.POST("/product/:id", controllers.UpdateProduct)
		api.GET("/products/category/:category_id", controllers.GetProductsByCategoryId) // New route

	}
}
