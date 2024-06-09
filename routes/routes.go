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
		api.GET("/products", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Products endpoint"})
		})
	}
}
