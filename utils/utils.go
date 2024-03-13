package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the structure for error responses.
type TypeErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
}

// SuccessResponse represents the structure for success responses.
type TypeSuccessResponse struct {
	Timestamp string      `json:"timestamp"`
	Status    int         `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
}

func ErrorResponse(c *gin.Context, status int, err error) {
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

func SuccessResponseWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().String(),
		"status":    http.StatusOK,
		"message":   message,
	})
}

func SuccessResponseWithDataAndMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().String(),
		"status":    http.StatusOK,
		"data":      data,
		"message":   message,
	})
}
