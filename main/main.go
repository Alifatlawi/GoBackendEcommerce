package main

import (
	"ecommercebackend/db"
	"ecommercebackend/routes"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, falling back to environment variables")
	}

	gin.SetMode(gin.ReleaseMode)
	db.InitDB()
	server := gin.Default()

	// Serve static files
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

	port := getPort()
	log.Printf("Starting server on port %s", port)
	err = server.Run(":" + port)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
