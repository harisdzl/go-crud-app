package customer_entity

import "github.com/harisquqo/quqo-challenge-1/domain/entity"


type Customer struct {
	entity.BaseModelWDelete
	ID uint64 `gorm:"primary_key;not null;" json:"id"`
	Name string `gorm:"size:255;not null;" json:"name"`
	Address string `gorm:"size:255;not null;" json:"address"`
	Latitude float64 `gorm:"type:numeric;not null;" json:"latitude"`
	Longitude float64 `gorm:"type:numeric;not null;" json:"longitude"`
}