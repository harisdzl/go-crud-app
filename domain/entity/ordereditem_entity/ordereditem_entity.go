package ordereditem_entity

import "github.com/harisquqo/quqo-challenge-1/domain/entity"

type OrderedItem struct {
	entity.BaseModelWDelete
	ID uint64 `json:"id"`
	OrderID int64 `gorm:"size:100;not null;" json:"order_id"`
	ProductID int64 `gorm:"size:255;not null;" json:"product_id"`
	Quantity int64 `gorm:"size:255;not null;" json:"quantity"` 
	UnitPrice float64 `gorm:"type:numeric;not null;" json:"unit_price"`
	TotalPrice float64 `gorm:"size:100;not null;" json:"total_price"`
}