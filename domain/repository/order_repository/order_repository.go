package order_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
)

type OrderRepository interface {
	SaveOrder(*order_entity.Order) (*order_entity.Order, map[string]string)
	GetOrder(int64) (*order_entity.Order, error)
	GetAllOrders() ([]order_entity.Order, error)
	UpdateOrder(*order_entity.Order) (*order_entity.Order, error)
	DeleteOrder(int64) error
}