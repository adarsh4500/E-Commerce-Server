package controllers

import (
	"Ecom/connections"
	"Ecom/models"
	"Ecom/postgres"
	"Ecom/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary View Orders
// @Description Retrieves a list of orders for the current user.
// @Tags Orders
// @Produce json
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /orders [get]
func ViewOrderHandler(c *gin.Context) {
	query := postgres.New(connections.DB)

	orders, err := query.ViewOrders(context.Background(), models.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, orders)
}


// @Summary Update Order Status
// @Description Updates the status of an order.
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body postgres.UpdateOrderStatusParams true "Order status information"
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /orders/updatestatus [post]
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

// @Summary View Order Items
// @Description Retrieves the items of a specific order by ID.
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID" format(uuid)
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /orders/{id} [get]
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
