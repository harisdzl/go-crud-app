package application

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
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
	c *gin.Context
}

func NewOrderApplication(p *base.Persistence, c *gin.Context) order_repository.OrderHandlerRepository {
	return &OrderApp{p, c}
}

func (a *OrderApp) CalculateTotalCost(rawOrder order_entity.RawOrder) float64 {
	span := a.p.Logger.Start(a.c, "application/CalculateTotalCost")
	defer span.End()
	var totalCost float64

	for productID, quantity := range rawOrder.Products {
		id, _ := strconv.ParseInt(productID, 10, 64)
		product, err := products.NewProductRepository(a.p, a.c).GetProduct(id); if err != nil {
			log.Println(err)
		}

		totalCost += (product.Price * float64(quantity))
	}	

	a.p.Logger.Info("application/CalculateTotalCost", map[string]interface{}{"total_cost": totalCost})
	return totalCost
}


func (a *OrderApp) SaveOrderFromRaw(rawOrder order_entity.RawOrder) (*order_entity.Order, error) {
	// Start a new span for the SaveOrderFromRaw function
	span := a.p.Logger.Start(a.c, "application/SaveOrderFromRaw", a.p.Logger.SetContextWithSpanFunc())
	defer span.End()
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
	
	a.p.Logger.Info("application/CalculateTotalCost", map[string]interface{}{"total_cost": totalCost})

	// Save the order
	repoOrder := orders.NewOrderRepository(a.p, a.c)
	savedOrder, err := repoOrder.SaveOrder(tx, &order)

	if err != nil {
		errTx = err
		return nil, errTx
	}

	repoOrderedItem := ordereditems.NewOrderedItemsRepository(a.p, a.c)
	for productID, quantity := range rawOrder.Products {
		productId, _ := strconv.ParseInt(productID, 10, 64)

		product, productErr := products.NewProductRepository(a.p, a.c).GetProduct(productId)
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

		inventoryRepo := inventories.NewInventoryRepository(a.p, a.c)
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
	span := a.p.Logger.Start(a.c, "application/GetOrder", a.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	repoOrder := orders.NewOrderRepository(a.p, a.c)
	return repoOrder.GetOrder(OrderId)
}

func (a *OrderApp) GetAllOrders() ([]order_entity.Order, error) {
	repoOrder := orders.NewOrderRepository(a.p, a.c)
	return repoOrder.GetAllOrders()
}
	
func (a *OrderApp) UpdateOrder(Order *order_entity.Order) (*order_entity.Order, error) {
	span := a.p.Logger.Start(a.c, "application/UpdateOrder", a.p.Logger.SetContextWithSpanFunc())
	defer span.End()
	repoOrder := orders.NewOrderRepository(a.p, a.c)
	return repoOrder.UpdateOrder(Order)
}

func (a *OrderApp) DeleteOrder(OrderId int64) error {
	repoOrder := orders.NewOrderRepository(a.p, a.c)
	return repoOrder.DeleteOrder(OrderId)
}


