package product_entity

import "github.com/harisquqo/quqo-challenge-1/domain/entity"

type Product struct {
	entity.BaseModelWDelete
	ID uint64 `json:"id"`
	Name string `gorm:"size:100;not null;" json:"name"`
	Description string `gorm:"size:255;not null;" json:"description"`
	Price float64 `gorm:"type:numeric;not null;" json:"price"`
	CategoryID uint64 `gorm:"size:100;not null;" json:"category_id"`
}

type ProductForInventory struct {
	entity.BaseModelWDelete
	ID uint64 `json:"id"`
	Name string `gorm:"size:100;not null;" json:"name"`
	Description string `gorm:"size:255;not null;" json:"description"`
	Price float64 `gorm:"type:numeric;not null;" json:"price"`
	CategoryID uint64 `gorm:"size:100;not null;" json:"category_id"`
	WarehouseID uint64 `gorm:"size:100;not null;" json:"warehouse_id"`
	Stock int `gorm:"size:255;not null;" json:"stock"`
}



