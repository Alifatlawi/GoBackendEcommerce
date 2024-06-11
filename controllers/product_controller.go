package controllers

import (
	"ecommercebackend/models"
	"ecommercebackend/repository"
	"fmt"
	"net/http"
	"os"
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
	// print request body
	fmt.Println(c.Request.Body)
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

	// Save the image to disk
	uploadPath := "uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		err := os.Mkdir(uploadPath, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
	}

	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExtension
	filePath := filepath.Join(uploadPath, newFileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	imgUrl := "/uploads/" + newFileName

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exists, Try another name"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "id": productID})
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		// Save the image to disk
		uploadPath := "uploads"
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			err := os.Mkdir(uploadPath, os.ModePerm)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
				return
			}
		}

		fileExtension := strings.ToLower(filepath.Ext(file.Filename))
		newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExtension
		filePath := filepath.Join(uploadPath, newFileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		// Update the product's ImgUrl field
		product.ImgUrl = filePath
	}

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
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	product, err := repository.GetProductById(product.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func GetProductsByCategoryId(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	products, err := repository.GetProductsByCategoryID(product.CategoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
