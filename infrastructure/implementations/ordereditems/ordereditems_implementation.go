package ordereditems

import (
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

type OrderedItemsRepo struct {
	p *base.Persistence
}

func NewOrderedItemsRepository(p *base.Persistence) *OrderedItemsRepo {
	return &OrderedItemsRepo{p}
}


func (o *OrderedItemsRepo) SaveOrderedItem(orderedItem *ordereditem_entity.OrderedItem) (*ordereditem_entity.OrderedItem, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", o.p)

	dbErr := map[string]string{}
	err := o.p.DB.Debug().Create(&orderedItem).Error
	if err != nil {
		fmt.Println("Failed to create orderedItem")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	cacheRepo.SetKey(fmt.Sprintf("%v_ORDEREDITEM", orderedItem.ID), orderedItem, time.Minute * 15)
	
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
