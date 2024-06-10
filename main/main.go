package main

import (
	"ecommercebackend/db"
	"ecommercebackend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	db.InitDB()
	server := gin.Default()

	// Apply the CORS middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Setup(server)
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
