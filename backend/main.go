package main

import (
	"log"
	"net/http"
	"radhaLLC/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Hardcoded admin credentials
var adminUsername = "admin"
var adminPassword = "password123"

func Login(c *gin.Context) {
	log.Println("Login endpoint hit")

	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
		})
		return
	}

	if loginDetails.Username == adminUsername && loginDetails.Password == adminPassword {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Login successful",
			"token":   "sample-token-123",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid credentials",
		})
	}
}

// Middleware to verify token
func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	log.Printf("Received token: %s", token) // Log received token

	// Check if token matches expected value
	if token != "Bearer sample-token-123" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access",
		})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true

	// Apply CORS middleware
	router.Use(cors.New(config))

	// Public routes
	router.POST("/login", Login)

	// Protected routes with authentication middleware
	protected := router.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("/employees", handlers.AddEmployee)
		protected.GET("/employees", handlers.GetEmployees)
		protected.GET("/employees/:id", handlers.GetEmployeeByID)
		protected.PUT("/employees/:id", handlers.UpdateEmployee)
		protected.DELETE("/employees/:id", handlers.DeleteEmployee)
	}

	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
