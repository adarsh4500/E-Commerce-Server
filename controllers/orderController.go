package controllers

import (
	"Ecom/connections"
	"Ecom/models"
	"Ecom/postgres"
	"Ecom/utils"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PlaceOrderHandler(c *gin.Context) {
	query := postgres.New(connections.DB)

	cart, err := query.ViewCart(context.Background(), models.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if len(cart) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("cart is empty"))
		return
	}

	oid, err := query.AddOrder(context.Background(), postgres.AddOrderParams{
		CustomerID:  models.UserID,
		TotalAmount: "0",
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	var total = float64(0.00)

	for _, item := range cart {
		i, err := query.AddOrderItem(context.Background(), postgres.AddOrderItemParams{
			OrderID:   oid,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			ID:        item.ProductID,
			Column5:   strconv.Itoa(int(item.Quantity)),
		})
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		f, err := strconv.ParseFloat(i, 64)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		total += f
	}

	strtotal := strconv.FormatFloat(total, 'f', -1, 64)
	orderdetails, err := query.UpdateOrderTotal(context.Background(), postgres.UpdateOrderTotalParams{
		ID:          oid,
		TotalAmount: strtotal,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	err = query.ClearCart(context.Background(), models.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, orderdetails)
}

func ViewOrderHandler(c *gin.Context) {
	query := postgres.New(connections.DB)

	orders, err := query.ViewOrders(context.Background(), models.UserID)
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

	query := postgres.New(connections.DB)

	orders, err := query.UpdateOrderStatus(context.Background(), param)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, orders)
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
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, items)
}
