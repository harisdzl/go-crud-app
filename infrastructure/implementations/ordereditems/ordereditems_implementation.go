package ordereditems

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

type OrderedItemsRepo struct {
	p *base.Persistence
	c *gin.Context
}

func NewOrderedItemsRepository(p *base.Persistence, c *gin.Context) *OrderedItemsRepo {
	return &OrderedItemsRepo{p, c}
}


func (o *OrderedItemsRepo) SaveOrderedItem(tx *gorm.DB, orderedItem *ordereditem_entity.OrderedItem) (*ordereditem_entity.OrderedItem, error) {
	channels := []string{"Zap", "Honeycomb"}
	loggerRepo, loggerErr := logger.NewLoggerRepository(channels, o.p, o.c, "implementations/SaveOrderedItem")
	defer loggerRepo.End()
	
	if loggerErr != nil {
		return nil, loggerErr
	}

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
		// tracer.Span.RecordError(err)
		fmt.Println("Failed to create orderedItem")
		fmt.Println(err)
		loggerRepo.Error(err.Error(), map[string]interface{}{})
		return nil, err
	}

	loggerRepo.Info("Ordered Item Saved", map[string]interface{}{"data": orderedItem})	
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
