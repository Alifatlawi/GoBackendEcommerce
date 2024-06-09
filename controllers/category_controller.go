package controllers

import (
	"ecommercebackend/models"
	"ecommercebackend/repository"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	categories, err := repository.GetAllCategories()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to query categories"})
		return
	}
	c.JSON(200, categories)
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	id, err := repository.CreateCategory(category)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: categories.name" {
			c.JSON(400, gin.H{"error": "Category already exists, Try another name"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(201, gin.H{"id": id})
}
