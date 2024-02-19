package ordereditem_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
)

type OrderedItemRepository interface {
	SaveOrderedItem(*ordereditem_entity.OrderedItem) (*ordereditem_entity.OrderedItem, map[string]string)
	SaveRawOrderItems(map[string]int64, int64) map[string]string
	GetAllOrderedItems() ([]ordereditem_entity.OrderedItem, error)
	GetAllOrderedItemsForOrder(int64) ([]ordereditem_entity.OrderedItem, error)
	ReverseOrder([]ordereditem_entity.OrderedItem) map[string]string
}