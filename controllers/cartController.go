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

func AddToCartHandler(c *gin.Context) {
	var cartItem models.CartItem

	err := c.ShouldBindJSON(&cartItem)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if cartItem.Quantity <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("quantity must be greater than 0"))
		return
	}

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

	product, err := query.GetProductById(context.Background(), cartItem.ProductID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, errors.New("encountered error while fetching product details"))
		return
	}

	if product.StockQuantity < cartItem.Quantity {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("not enough items in stock. requested:"+strconv.FormatInt(int64(cartItem.Quantity), 10)+", available"+strconv.FormatInt(int64(product.StockQuantity), 10)))
		return
	}

	item, err := query.AddToCart(context.Background(), postgres.AddToCartParams{
		UserID:    userID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, item)
}

func RemoveFromCartHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
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
	cart, err := query.RemoveFromCart(context.Background(), postgres.RemoveFromCartParams{
		UserID:    userID,
		ProductID: id,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, cart)
}

func ClearCartHandler(c *gin.Context) {
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
	err = query.ClearCart(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponseWithMessage(c, "Cart has been cleared")
}

func ViewCartHandler(c *gin.Context) {
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
	cart, err := query.ViewCart(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, cart)
}

func PlaceOrderHandler(c *gin.Context) {
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
	// Start transaction
	tx, err := connections.DB.Begin()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	txx := postgres.New(tx)
	cart, err := txx.ViewCart(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	if len(cart) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("cart is empty"))
		return
	}
	// Check stock for each item
	for _, item := range cart {
		product, err := txx.GetProductById(context.Background(), item.ProductID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		if product.StockQuantity < item.Quantity {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("not enough stock for product: "+product.Name))
			return
		}
	}
	// Place order
	orderID, err := txx.AddOrder(context.Background(), postgres.AddOrderParams{
		CustomerID:  userID,
		TotalAmount: "0.0",
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	total := 0.0
	for _, item := range cart {
		product, _ := txx.GetProductById(context.Background(), item.ProductID)
		price, _ := strconv.ParseFloat(product.Price, 64)
		subtotal := float64(item.Quantity) * price
		_, err := txx.AddOrderItem(context.Background(), postgres.AddOrderItemParams{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Column3:   item.Quantity,
		})
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		total += subtotal
		// Update stock
		_, err = txx.UpdateProductById(context.Background(), postgres.UpdateProductByIdParams{
			ID:            item.ProductID,
			Name:          product.Name,
			Price:         product.Price,
			Description:   product.Description,
			StockQuantity: product.StockQuantity - item.Quantity,
		})
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
	}
	orderdetails, err := txx.UpdateOrderTotal(context.Background(), postgres.UpdateOrderTotalParams{
		ID:          orderID,
		TotalAmount: strconv.FormatFloat(total, 'f', 2, 64),
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	err = txx.ClearCart(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, orderdetails)
}
