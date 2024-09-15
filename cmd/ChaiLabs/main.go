package main

import (
	"ChaiLabs/middleware"
	"ChaiLabs/repository"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file

	err := godotenv.Load(filepath.Join("E:/Go ChaiLabs backend", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server := gin.Default()
	//setting the base path for Api's end points
	basepath := server.Group("/api")
	server.Use(middleware.CORSMiddleware())

	// Connect to MongoDB
	repo, err := repository.NewMongoRepository()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Close MongoDB connection
	defer repo.Disconnect()

	// Register routes
	server.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Go ChaiLab backend!",
		})
	})

	//Adding base path to all services

	//Setting up auth service base path
	authController := repo.GetUserAuthController()
	authController.RegisterUserAuthRoutes(basepath)

	//Setting up user service base path
	userController := repo.GetUserController()
	userController.RegisterUserRoutes(authController.AuthService, basepath)

	// Run the server
	server.Run(":8000")
}
