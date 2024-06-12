package controllers

import (
	"ecommercebackend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteAllData(c *gin.Context) {
	err := repository.DeleteAllData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All data deleted successfully"})
}
