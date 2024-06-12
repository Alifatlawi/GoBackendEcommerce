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

		// Public routes
		api.GET("/categories", controllers.GetCategories)
		api.GET("/products", controllers.GetProducts)
		api.GET("/product/:id", controllers.GetProductById)
		api.GET("/products/category/:category_id", controllers.GetProductsByCategoryId)

		// Routes that require authentication
		api.Use(middleware.JWTAuth())
		{
			api.POST("/categories", controllers.CreateCategory)
			api.PUT("/categories", controllers.UpdateCategory)
			api.DELETE("/categories/:id", controllers.DeleteCategory)

			api.POST("/products", controllers.CreateProduct)
			api.PUT("/product/:id", controllers.UpdateProduct)
			api.PUT("/product/:id/image", controllers.UpdateProductImage)
			api.DELETE("/product/:id", controllers.DeleteProduct)

			// Route to delete all data
			api.DELETE("/delete-all", controllers.DeleteAllData)
		}
	}
}
