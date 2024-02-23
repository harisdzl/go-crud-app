package order_entity

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
)

type Order struct {
	entity.BaseModelWDelete
	ID uint64 `json:"id"`
	TotalCost float64 `gorm:"size:100;not null;" json:"total_cost"`
	TotalFees float64 `gorm:"size:255;not null;" json:"total_fees"`
	TotalCheckout float64 `gorm:"size:255;not null;" json:"total_checkout"` 
	CustomerID int64 `gorm:"not null;" json:"customer_id"`
	WarehouseID uint64 `gorm:"size:100;not null;" json:"warehouse_id"`
	Status string `gorm:"size:255;not null;" json:"status"`
}


type RawOrder struct {
	entity.BaseModelWDelete
	CustomerID int64 `gorm:"not null;" json:"customer_id"`
	WarehouseID uint64 `gorm:"size:100;not null;" json:"warehouse_id"`
	Status string `gorm:"size:255;not null;" json:"status"`
	Products  map[string]int64 `json:"products"`
}