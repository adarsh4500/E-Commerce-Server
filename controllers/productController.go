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

// @Summary Get All Products
// @Description Retrieves a list of all products.
// @Tags Products
// @Produce json
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /products [get]
func GetAllProductsHandler(c *gin.Context) {
	query := postgres.New(connections.DB)

	products, err := query.GetProducts(context.Background())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, products)
}

// @Summary Get Product by ID
// @Description Retrieves product details by ID.
// @Tags Products
// @Produce json
// @Param id path string true "Product ID" format(uuid)
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /products/{id} [get]
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

// @Summary Add New Product
// @Description Creates a new product.
// @Tags Products
// @Accept json
// @Produce json
// @Param request body models.Product true "New product information"
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /products/new [post]
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

// @Summary Update Product by ID
// @Description Updates product details by ID.
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" format(uuid)
// @Param request body models.EditProduct true "Fields to update"
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /products/update/{id} [post]
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

// @Summary Delete Product by ID
// @Description Deletes a product by ID.
// @Tags Products
// @Produce json
// @Param id path string true "Product ID" format(uuid)
// @Success 200 {object} utils.TypeSuccessResponse
// @Failure 400 {object} utils.TypeErrorResponse
// @Failure 500 {object} utils.TypeErrorResponse
// @Router /products/delete/{id} [post]
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
