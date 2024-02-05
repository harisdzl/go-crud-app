package application

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	order_repository "github.com/harisquqo/quqo-challenge-1/domain/repository/Order_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/orders"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type OrderApp struct {
	p *base.Persistence
}

func NewOrderApplication(p *base.Persistence) order_repository.OrderRepository {
	return &OrderApp{p}
}

type OrderAppInterface interface {
	SaveOrder(*order_entity.Order) (*order_entity.Order, map[string]string)
	GetOrder(int64) (*order_entity.Order, error)
	GetAllorders() ([]order_entity.Order, error)
	UpdateOrder(*order_entity.Order) (*order_entity.Order, error)
	DeleteOrder(int64) error
}

func (a *OrderApp) SaveOrder(Order *order_entity.Order) (*order_entity.Order, map[string]string) {
	repoOrder := orders.NewOrderRepository(a.p)
	return repoOrder.SaveOrder(Order)
}

func (a *OrderApp) GetOrder(OrderId int64) (*order_entity.Order, error) {
	repoOrder := orders.NewOrderRepository(a.p)
	return repoOrder.GetOrder(OrderId)
}

func (a *OrderApp) GetAllOrders() ([]order_entity.Order, error) {
	repoOrder := orders.NewOrderRepository(a.p)
	return repoOrder.GetAllOrders()
}
	
func (a *OrderApp) UpdateOrder(Order *order_entity.Order) (*order_entity.Order, error) {
	repoOrder := orders.NewOrderRepository(a.p)
	return repoOrder.UpdateOrder(Order)
}

func (a *OrderApp) DeleteOrder(OrderId int64) error {
	repoOrder := orders.NewOrderRepository(a.p)
	return repoOrder.DeleteOrder(OrderId)
}


