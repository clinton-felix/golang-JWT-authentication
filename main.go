package main

import (
	"log"
	"os"

	routes "github.com/clinton-felix/golang-JWT-auth-project/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main()  {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading the .env file..")
	}
	// set up port with .env value
	port := os.Getenv("PORT")
	if port == ""{
		port = "9000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Defining GET methods for API-1 nand API-2
	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success":"Access granted for api-1"})
	})

	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success":"Access granted for api-2"})
	})

	router.Run(":" + port)
}