package application

import (
	"fmt"
	"strconv"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/ordereditem_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/ordereditems"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/products"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type OrderedItemApp struct {
	p *base.Persistence
}

func NewOrderedItemApplication(p *base.Persistence) ordereditem_repository.OrderedItemRepository {
	return &OrderedItemApp{p}
}

func (a *OrderedItemApp) SaveOrderedItem(orderedItem *ordereditem_entity.OrderedItem) (*ordereditem_entity.OrderedItem, map[string]string) {
	repoOrderedItem := ordereditems.NewOrderedItemsRepository(a.p)
	return repoOrderedItem.SaveOrderedItem(orderedItem)
}

func (a *OrderedItemApp) SaveRawOrderItems(rawOrderProducts map[string]int64, orderID int64) map[string]string {
	for productID, quantity := range rawOrderProducts {
		productId, _ := strconv.ParseInt(productID, 10, 64)

		product, productErr := products.NewProductRepository(a.p).GetProduct(productId)
		if productErr != nil {
			return map[string]string{"error": "failed to retrieve product"}
		}


		orderedItem := ordereditem_entity.OrderedItem{
			OrderID:   orderID, // Assign the order ID to the ordered item
			ProductID: productId,
			Quantity:  quantity,
			UnitPrice: product.Price,
			TotalPrice: product.Price * float64(quantity), 
		}

		inventoryApp := NewInventoryApplication(a.p)
		inventoryApp.ReduceInventory(productId, quantity)

		// Save ordered item
		a.SaveOrderedItem(&orderedItem)
	}

	return nil
}

func (a *OrderedItemApp) ReverseOrder(orderedItems []ordereditem_entity.OrderedItem) map[string]string {
	errorMap := make(map[string]string)

	for _, orderedItem := range orderedItems {
		// Calculate the quantity to add to inventory by reversing the ordered quantity
		quantityToAdd := orderedItem.Quantity

		// Increase inventory for the product
		err := NewInventoryApplication(a.p).IncreaseInventory(orderedItem.ProductID, quantityToAdd)
		if err != nil {
			errorMap[fmt.Sprintf("product_%d", orderedItem.ProductID)] = err.Error()
		}
	}

	// Return error map, if any
	return errorMap
}

	func (a *OrderedItemApp) GetAllOrderedItems() ([]ordereditem_entity.OrderedItem, error) {
	repoOrderedItem := ordereditems.NewOrderedItemsRepository(a.p)
	return repoOrderedItem.GetAllOrderedItems()
}

func (a *OrderedItemApp) GetAllOrderedItemsForOrder(orderId int64) ([]ordereditem_entity.OrderedItem, error) {
	repoOrderedItem := ordereditems.NewOrderedItemsRepository(a.p)
	return repoOrderedItem.GetAllOrderedItemsForOrder(orderId)
}
