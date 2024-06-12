package controllers

import (
	"ecommercebackend/models"
	"ecommercebackend/repository"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetProducts Handle Get Products
func GetProducts(c *gin.Context) {
	products, err := repository.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// CreateProduct Handle Create Product
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Parse form data
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		fmt.Println("Binding Error:", err) // Debugging log
		return
	}

	// Debugging: Print the parsed product
	fmt.Printf("Parsed Product: %+v\n", product)

	// Get the category name from the form data
	categoryName := c.PostForm("category_name")

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// Generate a unique filename
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExtension

	// Upload the file to Azure
	imgUrl, err := uploadToAzure(file, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}
	product.ImgUrl = imgUrl

	// Check if the category name exists
	existingCategory, err := repository.GetCategoryByName(categoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query category"})
		return
	}
	if existingCategory.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category name"})
		return
	}

	// Set the category ID from the existing category
	product.CategoryId = existingCategory.ID

	// Create product in the repository
	productID, err := repository.CreateProduct(product)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: products.name") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exists, try another name"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "id": productID})
}

// UpdateProduct Handle Update Product
func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	var product models.Product

	// Handle different content types
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}
		if err := c.ShouldBind(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	product.ID = id
	if err := repository.UpdateProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// UpdateProductImage Handle Update Product Image
func UpdateProductImage(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	product, err := repository.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		return
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// Generate a unique filename
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExtension

	// Upload the file to Azure
	imgUrl, err := uploadToAzure(file, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Delete the old image from Azure storage
	if err := deleteFromAzure(product.ImgUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old image from Azure storage"})
		return
	}

	// Update the product's ImgUrl field
	product.ImgUrl = imgUrl

	if err := repository.UpdateProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product image updated successfully"})
}

// DeleteProduct Handle Delete Product
func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	product, err := repository.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		return
	}

	if product.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := repository.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	// Delete the image from Azure storage
	if err := deleteFromAzure(product.ImgUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image from Azure storage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product and associated image deleted successfully"})
}

// GetProductById Handle Get Product by ID
func GetProductById(c *gin.Context) {
	idParam := c.Param("id")
	id := idParam

	product, err := repository.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		return
	}

	if product.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductsByCategoryId Handle Get Products by Category ID
func GetProductsByCategoryId(c *gin.Context) {
	categoryIdParam := c.Param("category_id")
	categoryId := categoryIdParam

	products, err := repository.GetProductsByCategoryID(categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
