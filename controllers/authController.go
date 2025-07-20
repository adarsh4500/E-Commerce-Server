package controllers

import (
	"Ecom/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

// parseJWTToken is a helper function to parse and validate JWT tokens
func parseJWTToken(c *gin.Context) (jwt.MapClaims, error) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}

func Authenticate(c *gin.Context) {
	claims, err := parseJWTToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"timestamp": time.Now().String(),
			"status":    http.StatusUnauthorized,
			"message":   "authentication failed",
		})
		return
	}

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

	if float64(time.Now().Unix()) >= exp {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"timestamp": time.Now().String(),
			"status":    http.StatusUnauthorized,
			"message":   "token expired",
		})
		return
	}

	c.Set("user_id", idStr)
	c.Set("email", email)
	c.Next()
}

func RequireAdmin(c *gin.Context) {
	claims, err := parseJWTToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"timestamp": time.Now().String(),
			"status":    http.StatusUnauthorized,
			"message":   "authentication failed",
		})
		return
	}

	role, roleOk := claims["role"].(string)
	if !roleOk || role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"timestamp": time.Now().String(),
			"status":    http.StatusForbidden,
			"message":   "admin access required",
		})
		return
	}

	c.Next()
}
