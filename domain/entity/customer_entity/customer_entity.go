package customer_entity

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/utils/security"
	"gorm.io/gorm"
)


type Customer struct {
	entity.BaseModelWDelete
	ID int64 `gorm:"primary_key;not null;" json:"id"`
	Name string `gorm:"size:255;not null;" json:"name"`
	Address string `gorm:"size:255;not null;" json:"address"`
	Latitude float64 `gorm:"type:numeric;not null;" json:"latitude"`
	Longitude float64 `gorm:"type:numeric;not null;" json:"longitude"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}



//BeforeSave is a gorm hook
func (c *Customer) BeforeSave(tx *gorm.DB) error {
	hashPassword, err := security.Hash(c.Password)
	if err != nil {
		return err
	}
	c.Password = string(hashPassword)
	return nil
}