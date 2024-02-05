package category_entity

import "github.com/harisquqo/quqo-challenge-1/domain/entity"

type Category struct {
	entity.BaseModelWDelete
	ID uint64 `json:"id"`
	ParentID int64 `gorm:"size:100;not null;" json:"parent_id"`
	Name string `gorm:"size:100;not null;" json:"name"`
}
