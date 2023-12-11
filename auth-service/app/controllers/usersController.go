package controllers

import (
	"fmt"
	"music_streaming_service/auth-service/app/initializers"
	"music_streaming_service/auth-service/app/models"
	"net/http"
	"os"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SingUp(c *gin.Context) {
	// Get the email/passwd off req body
	var body struct {
		Email    string          `form:"Email"`
		Password string          `form:"Password"`
		Type     models.UserType `form:"Usertype"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Validate the email address using the isValidEmail function
	if !isValidEmail(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address",
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	if body.Type != models.Admin && body.Type != models.Editor && body.Type != models.Subscriber {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user type",
		})
		return
	}

	// Create user
	user := models.User{Email: body.Email, Password: string(hash), Type: body.Type}

	result := initializers.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	// Respond
	// c.JSON(http.StatusOK, gin.H{})
	c.Redirect(http.StatusFound, "/login")
}

func Login(c *gin.Context) {
	// Get the email/passwd off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Validate the email address using the isValidEmail function
	if !isValidEmail(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address",
		})
		return
	}

	// Look uo requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 || body.Email == os.Getenv("ADM_EMAIL") {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User was not found",
		})
		return
	}

	// Compare sent in pass with saved pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Generate a JWT token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subj":  user.ID,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
		"email": user.Email,
		"pwd":   user.Password,
		"type":  user.Type,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, os.Getenv("MUSIC_SERVICE_URL"))
}

func Validate(c *gin.Context) {
	// Get the token from the cookie
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Parse the JWT token
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Invalid signing method")
		}
		// Return the secret key used for signing
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user information from the claims
		userID := claims["subj"].(float64) // Assuming "subj" is the user ID
		userType := claims["type"].(string)
		userEmail := claims["email"].(string)

		// You can now use the extracted user information as needed
		c.JSON(http.StatusOK, gin.H{
			"id":    int(userID),
			"email": userEmail,
			"type":  userType,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
	}
}

func AdminLogin(c *gin.Context) {
	// Get the email/passwd off req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Validate the email address using the isValidEmail function
	if !isValidEmail(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address",
		})
		return
	}

	// Look uo requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 || body.Email != os.Getenv("ADM_EMAIL") {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User was not found",
		})
		return
	}

	// Compare sent in pass with saved pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Generate a JWT token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subj":  user.ID,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
		"email": user.Email,
		"pwd":   user.Password,
		"type":  user.Type,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, "/admin_panel")
	// c.JSON(http.StatusOK, gin.H{})
}

func isValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	// Compile the regular expression
	reg := regexp2.MustCompile(emailPattern, 0)
	// Use the MatchString method to check if the email matches the pattern
	isMatch, _ := reg.MatchString(email)
	return isMatch
}
