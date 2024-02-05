package image_entity

import "github.com/harisquqo/quqo-challenge-1/domain/entity"

type ImageRaw struct {
	entity.BaseModelWDelete
	ID uint64 `json:"id"`
	ProductID int64 `gorm:"size:100;not null;" json:"product_id"`
	Caption string `gorm:"not null;" json:"caption"`
}

type Image struct {
	entity.BaseModelWDelete
	ID uint64 `json:"id"`
	ProductID int64 `gorm:"size:100;not null;" json:"product_id"`
	Caption string `gorm:"not null;" json:"caption"`
	Url string `gorm:"not null;" json:"url"`
}
