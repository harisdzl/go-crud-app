package order_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	SaveOrder(*gorm.DB, *order_entity.Order) (*order_entity.Order, error)
	GetOrder(int64) (*order_entity.Order, error)
	GetAllOrders() ([]order_entity.Order, error)
	UpdateOrder(*order_entity.Order) (*order_entity.Order, error)
	DeleteOrder(int64) error
}



type OrderHandlerRepository interface {
	SaveOrderFromRaw(order_entity.RawOrder) (*order_entity.Order, error)
	GetOrder(int64) (*order_entity.Order, error)
	GetAllOrders() ([]order_entity.Order, error)
	UpdateOrder(*order_entity.Order) (*order_entity.Order, error)
	DeleteOrder(int64) error
}