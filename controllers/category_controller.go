package controllers

import (
	"ecommercebackend/models"
	"ecommercebackend/repository"
	"net/http"
	"strings"

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
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	categoryID, err := repository.CreateCategory(newCategory)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: categories.name") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category already exists, Try another name"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "id": categoryID})
}
