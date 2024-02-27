package ordereditems

import (
	"context"
	"errors"
	"fmt"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type OrderedItemsRepo struct {
	p *base.Persistence
	c context.Context
}

func NewOrderedItemsRepository(p *base.Persistence, c context.Context) *OrderedItemsRepo {
	return &OrderedItemsRepo{p, c}
}


func (o *OrderedItemsRepo) SaveOrderedItem(tx *gorm.DB, orderedItem *ordereditem_entity.OrderedItem) (*ordereditem_entity.OrderedItem, error) {
	tracer := otel.Tracer("implementations.inventories.ReduceInventory")
	_, span := tracer.Start(o.c, "implementations.inventories.ReduceInventory")
	defer span.End()

	if tx == nil {
		var errTx error
		tx := o.p.DB.Begin()
		if tx.Error != nil {
			return nil, errors.New("failed to start transaction")
		}
	
		// Defer rollback in case of panic
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			} else if errTx != nil {
				tx.Rollback()
			} else {
				errC := tx.Commit().Error
				if errC != nil {
					tx.Rollback()
				}
			}
		}()
	}

	err := tx.Debug().Create(&orderedItem).Error
	if err != nil {
		span.RecordError(err)
		fmt.Println("Failed to create orderedItem")
		fmt.Println(err)
		return nil, err
	}

	
	return orderedItem, nil
}


func (o *OrderedItemsRepo) GetAllOrderedItems() ([]ordereditem_entity.OrderedItem, error) {
	var orderedItems []ordereditem_entity.OrderedItem
	err := o.p.DB.Debug().Find(&orderedItems).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return orderedItems, nil
}

func (o *OrderedItemsRepo) GetAllOrderedItemsForOrder(orderId int64) ([]ordereditem_entity.OrderedItem, error) {
	var orderedItems []ordereditem_entity.OrderedItem

	err := o.p.DB.Debug().Where("order_id = ?", orderId).Find(&orderedItems).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return orderedItems, nil
}
