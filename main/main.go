package main

import (
	"ecommercebackend/db"
	"ecommercebackend/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	db.InitDB()
	server := gin.Default()
	routes.Setup(server)
	err := server.Run(":8080")
	if err != nil {
		return
	}
}
