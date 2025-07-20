package controllers

import (
	"Ecom/connections"
	"Ecom/postgres"
	"Ecom/utils"
	"context"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ViewOrderHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("user not found in context"))
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid user id"))
		return
	}
	query := postgres.New(connections.DB)
	orders, err := query.ViewOrders(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, orders)
}

func UpdateOrderStatusHandler(c *gin.Context) {
	var param postgres.UpdateOrderStatusParams
	err := c.BindJSON(&param)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	allowed := map[string]bool{"Pending": true, "Shipped": true, "Delivered": true, "Cancelled": true, "Processing": true}
	if !allowed[param.Status] {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid order status"))
		return
	}
	query := postgres.New(connections.DB)
	order, err := query.UpdateOrderStatus(context.Background(), param)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	// Check authorization
	claims, _ := parseJWTToken(c)
	role, _ := claims["role"].(string)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("user not found in context"))
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid user id"))
		return
	}
	if role != "admin" && order.CustomerID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, errors.New("not authorized to update this order"))
		return
	}
	utils.SuccessResponse(c, order)
}

func ViewOrderItemsHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	query := postgres.New(connections.DB)
	items, err := query.ViewOrderItems(context.Background(), id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	// Check if order belongs to user
	if len(items) > 0 {
		orderID := items[0].OrderID
		order, err := query.GetOrderById(context.Background(), orderID)
		if err == nil {
			userIDStr, exists := c.Get("user_id")
			if exists {
				userID, _ := uuid.Parse(userIDStr.(string))
				if order.CustomerID != userID {
					utils.ErrorResponse(c, http.StatusForbidden, errors.New("not authorized to view this order's items"))
					return
				}
			}
		}
	}
	utils.SuccessResponse(c, items)
}

func OrderHistoryHandler(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("user not found in context"))
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid user id"))
		return
	}
	query := postgres.New(connections.DB)
	orders, err := query.ViewOrders(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	var result []map[string]interface{}
	for _, order := range orders {
		items, _ := query.ViewOrderItems(context.Background(), order.ID)
		result = append(result, map[string]interface{}{
			"id":           order.ID,
			"order_date":   order.OrderDate,
			"total_amount": order.TotalAmount,
			"status":       order.Status,
			"items":        items,
		})
	}
	utils.SuccessResponse(c, result)
}

// Get all pending orders
func AdminAllOrdersHandler(c *gin.Context) {
	query := postgres.New(connections.DB)
	orders, err := query.GetAllPendingOrders(context.Background())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	var result []map[string]interface{}
	for _, order := range orders {
		items, _ := query.ViewOrderItems(context.Background(), order.ID)
		result = append(result, map[string]interface{}{
			"id":           order.ID,
			"customer_id":  order.CustomerID,
			"order_date":   order.OrderDate,
			"total_amount": order.TotalAmount,
			"status":       order.Status,
			"items":        items,
		})
	}
	utils.SuccessResponse(c, result)
}
