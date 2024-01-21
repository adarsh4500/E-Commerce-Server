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


// @Summary Add to Cart
// @Description Adds a product to the user's cart.
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body models.CartItem true "Cart item information"
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /cart/new [post]
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


// @Summary Remove from Cart
// @Description Removes a product from the user's cart.
// @Tags Cart
// @Produce json
// @Param id path string true "Product ID" format(uuid)
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /cart/remove/{id} [post]
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

// @Summary Clear Cart
// @Description Clears all items from the user's cart.
// @Tags Cart
// @Produce json
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /cart/clear [post]
func ClearCartHandler(c *gin.Context){
	query := postgres.New(connections.DB)

	err:= query.ClearCart(context.Background(),models.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, gin.H{ "message":"cart has been cleared"})
}

// @Summary View Cart
// @Description Retrieves the items in the user's cart.
// @Tags Cart
// @Produce json
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /cart [get]
func ViewCartHandler(c *gin.Context){
	query := postgres.New(connections.DB)

	cart, err:= query.ViewCart(context.Background(),models.UserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, cart)
}

// @Summary Place Order
// @Description Places an order using the items in the user's cart.
// @Tags Cart
// @Produce json
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /cart/checkout [post]
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