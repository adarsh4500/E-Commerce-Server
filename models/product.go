package models

type Product struct {
	Name          string `json:"name" binding:"required"`
	Price         string `json:"price" binding:"required"`
	Description   string `json:"description"`
	StockQuantity int32  `json:"stock_quantity" binding:"required"`
}

type EditProduct struct {
	Name          *string `json:"name"`
	Price         *string `json:"price"`
	Description   *string `json:"description"`
	StockQuantity *int32    `json:"stock_quantity"`
}
