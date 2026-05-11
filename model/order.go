package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uint64     `json:"order_id"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	LineItems   []LineItem `json:"line_items"`
	CreatedAt   *time.Time `json:"created_time"`
	ShippedAt   *time.Time `json:"shipped_time"`
	CompletedAt *time.Time `json:"completed_time"`
}

type LineItem struct {
	ItemID   uuid.UUID `json:"item_id"`
	Quantity uint16    `json:"quantity"`
	Price    float32   `json:"price"`
}
