package controllers

import (
	"Ecom/config"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Authenticate(c *gin.Context) {
	// Get Token off Cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"timestamp": time.Now().String(),
			"status":    http.StatusUnauthorized,
			"message":   err.Error(),
		})
		return
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"timestamp": time.Now().String(),
				"status":    http.StatusUnauthorized,
				"message":   "error occurred while parsing token",
			})
			return nil, nil
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"timestamp": time.Now().String(),
			"status":    http.StatusUnauthorized,
			"message":   err.Error(),
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, emailOk := claims["email"].(string)
		idStr, idOk := claims["id"].(string)
		exp, expOk := claims["exp"].(float64)
		if !emailOk || !idOk || !expOk {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"timestamp": time.Now().String(),
				"status":    http.StatusUnauthorized,
				"message":   "invalid jwt claims",
			})
			return
		}
		if float64(time.Now().Unix()) < exp {
			c.Set("user_id", idStr)
			c.Set("email", email)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"timestamp": time.Now().String(),
				"status":    http.StatusUnauthorized,
				"message":   "token expired",
			})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"timestamp": time.Now().String(),
			"status":    http.StatusUnauthorized,
			"message":   "invalid jwt token",
		})
		return
	}
}
