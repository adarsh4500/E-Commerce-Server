package models

type Product struct {
	Name          string  `json:"name" binding:"required,min=2,max=255"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	Description   string  `json:"description"`
	StockQuantity int32   `json:"stock_quantity" binding:"required,gte=0"`
}

type EditProduct struct {
	Name          *string  `json:"name"`
	Price         *float64 `json:"price"`
	Description   *string  `json:"description"`
	StockQuantity *int32   `json:"stock_quantity"`
}
