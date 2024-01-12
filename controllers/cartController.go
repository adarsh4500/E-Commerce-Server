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
	_ "github.com/lib/pq"
)

func AddToCartHandler(c *gin.Context) {
	var cartItem models.CartItem

	err := c.ShouldBindJSON(&cartItem)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	query := postgres.New(connections.DB)

	product, err := query.GetProductById(context.Background(), cartItem.ProductID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, errors.New("encoutered error while fetching product details"))
		return
	}

	if product.StockQuantity < cartItem.Quantity {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("not enough items in stock. requested:"+strconv.FormatInt(int64(cartItem.Quantity), 10)+", available"+strconv.FormatInt(int64(product.StockQuantity), 10)))
		return
	}

	item, err := query.AddToCart(context.Background(), postgres.AddToCartParams{
		UserID:    models.UserID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, item)
}

func RemoveFromCartHandler(c *gin.Context){
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	query := postgres.New(connections.DB)
	cart, err:= query.RemoveFromCart(c,postgres.RemoveFromCartParams{
		UserID: models.UserID,
		ProductID: id,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, cart)
}

func ClearCartHandler(c *gin.Context){
	query := postgres.New(connections.DB)

	err:= query.ClearCart(context.Background(),models.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, gin.H{ "message":"cart has been cleared"})
}

func ViewCartHandler(c *gin.Context){
	query := postgres.New(connections.DB)

	cart, err:= query.ViewCart(context.Background(),models.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, cart)
}