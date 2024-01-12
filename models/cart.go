package models

import "github.com/google/uuid"

var Cart []CartItem

type CartItem struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int32     `json:"quantity" binding:"required"`
}
