package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"music_streaming_service/auth-service/app/controllers"
	"music_streaming_service/auth-service/app/initializers"
	"music_streaming_service/auth-service/app/middleware"
	"music_streaming_service/auth-service/app/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
}

func main() {
	router := gin.Default()
	router.POST("/registration", controllers.SingUp)
	router.POST("/signin", controllers.Login)
	router.POST("/admin_signin", controllers.AdminLogin)
	router.GET("/valid", middleware.RequireAuth1, controllers.Validate)
	router.LoadHTMLGlob("resources/templates/*.html")
	router.Static("/static", "./resources/static")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Join LAD",
		})
	})

	router.GET("/admin_panel", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_panel.html", gin.H{
			"title": "LAD admin panel",
		})
	})

	router.GET("/login", middleware.RequireAuthCommon, func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Sign in to LAD",
		})
	})

	router.GET("/admin", middleware.RequireAuthAdmin, func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{
			"title": "Login to admin panel",
		})
	})

	// List all users
	router.GET("/admin/users", func(c *gin.Context) {
		var users []models.User
		initializers.DB.Find(&users)
		c.JSON(200, users)
	})

	// Update a user
	router.PUT("/admin/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Implement logic to update the user in the database (Gorm)
		var user models.User
		if err := initializers.DB.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		// Decode the JSON request body into a struct (you may need to create a UserUpdate struct)
		var updatedUserData models.User
		if err := c.BindJSON(&updatedUserData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// Update email and UserType
		user.Email = updatedUserData.Email
		user.Type = updatedUserData.Type

		if err := initializers.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	// Delete a user
	router.DELETE("/admin/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		initializers.DB.Delete(&models.User{}, id)
		c.JSON(200, gin.H{"message": "User deleted"})
	})

	// Start the Gin server in a separate goroutine
	go func() {
		err := router.Run()
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Start the SendPostReq in a separate goroutine with a 5-second delay
	go func() {
		time.Sleep(5 * time.Second)
		err := SendPostReq()
		if err != nil {
			fmt.Println("Error sending POST request:", err)
		}
	}()

	// Keep the main goroutine running to allow other goroutines to execute
	select {}
}

func SendPostReq() error {
	// Define the request URL
	port := os.Getenv("PORT")
	url := "http://localhost:" + port + "/registration"
	method := "POST"

	// Create the JSON payload
	payload := map[string]interface{}{
		"email":    os.Getenv("ADM_EMAIL"),
		"password": os.Getenv("ADM_PWD"),
		"type":     "admin",
	}

	// Marshal the payload into JSON format
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create a POST request with the JSON payload
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	// Set the request headers (optional)
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
}
