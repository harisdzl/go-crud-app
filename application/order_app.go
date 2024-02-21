package application

import (
	"errors"
	"log"
	"strconv"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/ordereditem_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/order_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/ordereditems"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/orders"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/products"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type OrderApp struct {
	p *base.Persistence
}

func NewOrderApplication(p *base.Persistence) order_repository.OrderHandlerRepository {
	return &OrderApp{p}
}

func (a *OrderApp) CalculateTotalCost(rawOrder order_entity.RawOrder) float64 {
	var totalCost float64
	for productID, quantity := range rawOrder.Products {
		id, _ := strconv.ParseInt(productID, 10, 64)
		product, err := products.NewProductRepository(a.p).GetProduct(id); if err != nil {
			log.Println(err)
		}

		totalCost += (product.Price * float64(quantity))
	}

	return totalCost
}


func (a *OrderApp) SaveOrderFromRaw(rawOrder order_entity.RawOrder) (*order_entity.Order, error) {
	var errTx error

	tx := a.p.DB.Begin()
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
	
	// Create an order entity
	order := order_entity.Order{
		CustomerID:  rawOrder.CustomerID,
		WarehouseID: rawOrder.WarehouseID,
		Status:      rawOrder.Status,
		TotalFees:   0,
	}

	// Calculates total costs of all the products
	totalCost := a.CalculateTotalCost(rawOrder)
	// Set other fields of the order entity
	order.TotalCost = totalCost
	totalCheckout := totalCost + order.TotalFees
	order.TotalCheckout = totalCheckout

	// Save the order
	repoOrder := orders.NewOrderRepository(a.p)
	savedOrder, err := repoOrder.SaveOrder(tx, &order)

	if err != nil {
		errTx = err
        return nil, errTx
    }

	repoOrderedItem := ordereditems.NewOrderedItemsRepository(a.p)
	for productID, quantity := range rawOrder.Products {
		productId, _ := strconv.ParseInt(productID, 10, 64)

		product, productErr := products.NewProductRepository(a.p).GetProduct(productId)
		if productErr != nil {
			errTx = err
			return nil, productErr
		}

		orderedItem := ordereditem_entity.OrderedItem{
			OrderID:    int64(order.ID), // Assign the order ID to the ordered item
			ProductID:  productId,
			Quantity:   quantity,
			UnitPrice:  product.Price,
			TotalPrice: product.Price * float64(quantity),
		}

		inventoryRepo := inventories.NewInventoryRepository(a.p)
		reduceInventoryErr := inventoryRepo.ReduceInventory(tx, productId, quantity)

		if reduceInventoryErr != nil {
			errTx = reduceInventoryErr
			return nil, errTx
		}
		// Save ordered item
		_, err := repoOrderedItem.SaveOrderedItem(tx, &orderedItem)

		if err != nil {
			errTx = err
			return nil, errTx
		}
	}

	return savedOrder, nil
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


