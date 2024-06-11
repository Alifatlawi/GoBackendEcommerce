package main

import (
	"ecommercebackend/db"
	"ecommercebackend/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	db.InitDB()
	server := gin.Default()
	// serve static files
	server.Static("/uploads", "./uploads")

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
	err := server.Run(":" + getPort())
	if err != nil {
		return
	}
}

func getPort() string {
	port := "8080"
	if value, exists := os.LookupEnv("PORT"); exists {
		port = value
	}
	return port
}
