package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, status int,  err error) {
	c.JSON(status, gin.H{
		"timestamp": time.Now().String(),
		"status":    status,
		"message":   err.Error(),
	})
}

func SuccessResponse(c *gin.Context, res interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().String(),
		"status":    http.StatusOK,
		"data":      res,
	})
}
