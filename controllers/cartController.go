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
	rows, err := query.ViewCart(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	var cart []models.CartItemWithProduct
	for _, row := range rows {
		cart = append(cart, models.CartItemWithProduct{
			ID:         row.ID,
			UserID:     row.UserID,
			ProductID:  row.ProductID,
			Quantity:   row.Quantity,
			ModifiedAt: row.ModifiedAt.Format("2006-01-02T15:04:05Z07:00"),
			Product: models.ProductDetails{
				ID:            row.ProductID_2,
				Name:          row.ProductName,
				Price:         row.ProductPrice,
				Description:   row.ProductDescription,
				StockQuantity: row.ProductStockQuantity,
			},
		})
	}
	utils.SuccessResponse(c, cart)
}

func CartCountHandler(c *gin.Context) {
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
	rows, err := query.ViewCart(context.Background(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	total := 0
	for _, row := range rows {
		total += int(row.Quantity)
	}
	utils.SuccessResponse(c, map[string]int{"count": total})
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
		// Reduce stock
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

func UpdateCartItemHandler(c *gin.Context) {
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
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid product id"))
		return
	}
	var body struct {
		Quantity int32 `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	query := postgres.New(connections.DB)
	if body.Quantity <= 0 {
		// Remove item if quantity is 0 or less
		_, err := query.RemoveFromCart(context.Background(), postgres.RemoveFromCartParams{
			UserID:    userID,
			ProductID: productID,
		})
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		utils.SuccessResponseWithMessage(c, "Item removed from cart")
		return
	}
	item, err := query.UpdateCartItem(context.Background(), postgres.UpdateCartItemParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  body.Quantity,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, item)
}
