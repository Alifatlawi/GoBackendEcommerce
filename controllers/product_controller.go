package controllers

import (
	"ecommercebackend/models"
	"ecommercebackend/repository"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	products, err := repository.GetAllProducts()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to query products"})
		return
	}
	c.JSON(200, products)
}

func CreateProduct(c *gin.Context) {
	// Parse form data
	name := c.PostForm("name")
	description := c.PostForm("description")
	categoryId, _ := strconv.Atoi(c.PostForm("category_id"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

	// Handle image upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed"})
		return
	}

	// Generate a unique filename
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExtension

	// Upload the file to AWS S3
	imgUrl, err := uploadToAzure(file, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Create product
	var product = models.Product{
		Name:        name,
		Description: description,
		CategoryId:  categoryId,
		ImgUrl:      imgUrl,
		Price:       price,
	}
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

func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Check content type to handle different payloads
	contentType := c.GetHeader("Content-Type")
	var product models.Product

	if strings.HasPrefix(contentType, "application/json") {
		// Handle JSON payload
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		// Extract fields from the form
		product.Id = id
		product.Name = c.PostForm("name")
		product.Description = c.PostForm("description")
		categoryId, err := strconv.Atoi(c.PostForm("category_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		product.CategoryId = categoryId
		price, err := strconv.ParseFloat(c.PostForm("price"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
			return
		}
		product.Price = price

		// Handle image upload
		file, err := c.FormFile("image")
		if err == nil {
			// Generate a unique filename
			fileExtension := strings.ToLower(filepath.Ext(file.Filename))
			newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExtension

			// Upload the file to Azure (adjust function if needed)
			imgUrl, err := uploadToAzure(file, newFileName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
				return
			}

			// Update the product's ImgUrl field
			product.ImgUrl = imgUrl
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	product.Id = id // Ensure the ID is set

	err = repository.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func DeleteProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	err := repository.DeleteProduct(product.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func GetProductById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := repository.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		return
	}

	if product.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func GetProductsByCategoryId(c *gin.Context) {
	categoryIdParam := c.Param("category_id")
	categoryId, err := strconv.Atoi(categoryIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	products, err := repository.GetProductsByCategoryID(categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
