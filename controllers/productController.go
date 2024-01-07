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
	_ "github.com/lib/pq"
)

func GetAllProductsHandler(c *gin.Context) {
	query := postgres.New(connections.DB)

	products, err := query.GetProducts(context.Background())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, products)
}

func GetProductHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	query := postgres.New(connections.DB)
	product, err := query.GetProductById(context.Background(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, product)
}

func NewProductHandler(c *gin.Context) {

	var newProduct models.Product
	err := c.ShouldBindJSON(&newProduct)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var param = postgres.AddProductParams{
		Name:          newProduct.Name,
		Description:   newProduct.Description,
		Price:         newProduct.Price,
		StockQuantity: newProduct.StockQuantity,
	}

	query := postgres.New(connections.DB)

	product, err := query.AddProduct(context.Background(), param)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, product)

}
func UpdateProductHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var updateFields models.EditProduct
	err = c.ShouldBindJSON(&updateFields)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	query := postgres.New(connections.DB)
	product, err := query.GetProductById(context.Background(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var param postgres.UpdateProductByIdParams
	param.ID = id
	param.Name = product.Name
	param.Description = product.Description
	param.Price = product.Price
	param.StockQuantity = product.StockQuantity

	if updateFields.Name != nil {
		param.Name = *updateFields.Name
	}
	if updateFields.Description != nil {
		param.Description = *updateFields.Description
	}
	if updateFields.Price != nil {
		param.Price = *updateFields.Price
	}
	if updateFields.StockQuantity != nil {
		param.StockQuantity = *updateFields.StockQuantity
	}

	update, err := query.UpdateProductById(context.Background(), param)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, update)
}

func DeleteProductHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	query := postgres.New(connections.DB)
	product, err := query.DeleteProductById(context.Background(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, product)
}
