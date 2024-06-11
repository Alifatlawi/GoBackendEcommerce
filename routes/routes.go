package routes

import (
	"ecommercebackend/controllers"
	"ecommercebackend/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)

		//public routes
		api.GET("/categories", controllers.GetCategories)
		api.GET("/products", controllers.GetProducts)
		api.GET("/product/:id", controllers.GetProductById)
		api.GET("/products/category/:category_id", controllers.GetProductsByCategoryId)

		// Routes that require authentication
		api.Use(middleware.JWTAuth())
		{
			api.POST("/categories", controllers.CreateCategory)
			api.PUT("/categories", controllers.UpdateCategory)    // Ensure UpdateCategory is defined
			api.DELETE("/categories", controllers.DeleteCategory) // Ensure DeleteCategory is defined
			api.POST("/products", controllers.CreateProduct)
			api.PUT("/products", controllers.UpdateProduct)
			api.DELETE("/products", controllers.DeleteProduct)
			api.POST("/product/:id", controllers.UpdateProduct)
			api.PUT("/product/:id", controllers.UpdateProduct)
		}
	}
}
