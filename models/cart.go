package models

import "github.com/google/uuid"

var Cart []CartItem

type CartItem struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int32     `json:"quantity" binding:"required"`
}

type CartItemWithProduct struct {
	ID         uuid.UUID      `json:"id"`
	UserID     uuid.UUID      `json:"user_id"`
	ProductID  uuid.UUID      `json:"product_id"`
	Quantity   int32          `json:"quantity"`
	ModifiedAt string         `json:"modified_at"`
	Product    ProductDetails `json:"product"`
}

type ProductDetails struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Price         string    `json:"price"`
	Description   string    `json:"description"`
	StockQuantity int32     `json:"stock_quantity"`
}
