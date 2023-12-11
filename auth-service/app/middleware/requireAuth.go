package middleware

import (
	"errors"
	"fmt"
	"music_streaming_service/auth-service/app/initializers"
	"music_streaming_service/auth-service/app/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth1(c *gin.Context) {
	// Get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/validate it

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach to reqk
		c.Set("user", user)

		// Continue
		c.Next()

	} else {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

func RequireAuthCommon(c *gin.Context) {
	// Get the cookie off req
	tokenString, err := c.Cookie("Authorization")
	var localError error = nil

	if err != nil {
		// localError = errors.New("There is no cookies!")
		// fmt.Println(localError)
		err = nil
		c.Next()
	}

	// Decode/validate it

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			localError = errors.New("Experation time is over!")
			fmt.Println(localError)
		}

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 || claims["type"] == "admin" {
			localError = errors.New("User was not found!")
			fmt.Println(localError)
		}

		// Attach to reqk
		c.Set("user", user)

		// Continue
		if localError != nil {
			c.Next()
		} else {
			c.Redirect(http.StatusFound, "/music_service")
		}

	} else {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

func RequireAuthAdmin(c *gin.Context) {
	// Get the cookie off req
	tokenString, err := c.Cookie("Authorization")
	var localError error = nil

	if err != nil {
		// localError = errors.New("There is no cookies!")
		// fmt.Println(localError)
		err = nil
		c.Next()
	}

	// Decode/validate it

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			localError = errors.New("Experation time is over!")
			fmt.Println(localError)
		}

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 || claims["type"] != "admin" {
			localError = errors.New("User was not found!")
			fmt.Println(localError)
		}

		// Attach to reqk
		c.Set("user", user)

		// Continue
		if localError != nil {
			c.Next()
		} else {
			c.Redirect(http.StatusFound, "/admin_panel")
		}

	} else {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
