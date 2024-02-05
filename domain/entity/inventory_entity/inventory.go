package inventory_entity

import "github.com/harisquqo/quqo-challenge-1/domain/entity"

type Inventory struct {
	entity.BaseModelWDelete
	ProductID uint64 `gorm:"primary_key;not null;" json:"product_id"`
	WarehouseID uint64 `gorm:"size:100;not null;" json:"warehouse_id"`
	Stock int `gorm:"size:255;not null;" json:"stock"`
}

