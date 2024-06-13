package controllers

import (
	"ecommercebackend/models"
	"ecommercebackend/repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse sends a JSON response with the given status code and error message
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func GetCategories(c *gin.Context) {
	categories, err := repository.GetAllCategories()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query categories")
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	categoryID, err := repository.CreateCategory(newCategory)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: categories.name") {
			ErrorResponse(c, http.StatusBadRequest, "Category already exists, try another name")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create category")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "id": categoryID})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var json struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	category, err := repository.GetCategoryById(id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query category")
		return
	}
	if category.ID == "" {
		ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	category.Name = json.Name

	if err := repository.UpdateCategory(category); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: categories.name") {
			ErrorResponse(c, http.StatusBadRequest, "Category name already exists")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update category")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Delete all products that belong to this category
	if err := repository.DeleteProductsByCategoryId(id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete products by category")
		return
	}

	// Delete the category
	if err := repository.DeleteCategory(id); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category and associated products deleted successfully"})
}
