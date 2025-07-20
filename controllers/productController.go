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

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllProductsHandler(c *gin.Context) {
	query := postgres.New(connections.DB)

	search := c.Query("search")
	sort := c.Query("sort")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page := 1
	limit := 16
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	products, err := query.GetProductsSorted(context.Background(), postgres.GetProductsSortedParams{
		Column1: sql.NullString{String: search, Valid: search != ""},
		Column2: sort,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
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
		if err.Error() == "sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
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
	// Input validation
	if len(newProduct.Name) < 2 || len(newProduct.Name) > 255 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("name must be 2-255 characters"))
		return
	}
	if newProduct.Price <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("price must be greater than 0"))
		return
	}
	if newProduct.StockQuantity < 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, errors.New("stock_quantity must be >= 0"))
		return
	}

	var param = postgres.AddProductParams{
		Name:          newProduct.Name,
		Description:   newProduct.Description,
		Price:         strconv.FormatFloat(newProduct.Price, 'f', 2, 64),
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
		if err.Error() == "sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
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
		if len(*updateFields.Name) < 2 || len(*updateFields.Name) > 255 {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("name must be 2-255 characters"))
			return
		}
		param.Name = *updateFields.Name
	}
	if updateFields.Description != nil {
		param.Description = *updateFields.Description
	}
	if updateFields.Price != nil {
		if *updateFields.Price <= 0 {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("price must be greater than 0"))
			return
		}
		param.Price = strconv.FormatFloat(*updateFields.Price, 'f', 2, 64)
	}
	if updateFields.StockQuantity != nil {
		if *updateFields.StockQuantity < 0 {
			utils.ErrorResponse(c, http.StatusBadRequest, errors.New("stock_quantity must be >= 0"))
			return
		}
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
		if err.Error() == "sql: no rows in result set" {
			utils.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, product)
}

func GetProductsCountHandler(c *gin.Context) {
	search := c.Query("search")
	query := postgres.New(connections.DB)
	count, err := query.GetProductsCount(context.Background(), sql.NullString{String: search, Valid: search != ""})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, count)
}
