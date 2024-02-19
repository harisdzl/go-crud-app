package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/order_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

type OrderRepo struct {
	p *base.Persistence
}

func NewOrderRepository(p *base.Persistence) *OrderRepo {
	return &OrderRepo{p}
}

var _ order_repository.OrderRepository = &OrderRepo{}

func (o *OrderRepo) SaveOrder(order *order_entity.Order) (*order_entity.Order, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", o.p)

	dbErr := map[string]string{}
	err := o.p.DB.Debug().Create(&order).Error
	if err != nil {
		fmt.Println("Failed to create order")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	cacheRepo.SetKey(fmt.Sprintf("%v_ORDER", order.ID), order, time.Minute * 15)
	
	return order, nil
}


func (o *OrderRepo) GetOrder(id int64) (*order_entity.Order, error) {
	var order *order_entity.Order

	cacheRepo := cache.NewCacheRepository("Redis", o.p)
	_ = cacheRepo.GetKey(fmt.Sprintf("%v_ORDER", id), &order)
	if order == nil {
		err := o.p.DB.Debug().Where("id = ?", id).Take(&order).Error
		if err != nil {
			fmt.Println("Failed to get order")
		}
		if order != nil && order.ID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprintf("%v_ORDER", id), order, time.Minute * 15)
		}
	}


	return order, nil
}

func (o *OrderRepo) GetAllOrders() ([]order_entity.Order, error) {
	var orders []order_entity.Order
	err := o.p.DB.Debug().Find(&orders).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return orders, nil
}

func (o *OrderRepo) UpdateOrder(order *order_entity.Order) (*order_entity.Order, error) {
	cacheRepo := cache.NewCacheRepository("Redis", o.p)


	err := o.p.DB.Debug().Where("id = ?", order.ID).Updates(&order).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = cacheRepo.SetKey(fmt.Sprintf("%v_ORDER", order.ID), order, time.Minute * 15)

	return order, nil
}

func (o *OrderRepo) DeleteOrder(id int64) error {
	var order order_entity.Order

	err := o.p.DB.Debug().Where("id = ?", id).Delete(&order).Error
	
	cacheRepo := cache.NewCacheRepository("Redis", o.p)

	cacheRepo.DelKey(fmt.Sprintf("%v_ORDER", id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}