package application

import (
	"log"
	"strconv"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/order_repository"
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


func (a *OrderApp) SaveOrderFromRaw(rawOrder order_entity.RawOrder) (*order_entity.Order, map[string]string) {
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
	savedOrder, err := repoOrder.SaveOrder(&order)

	// Call orderedItem application function
	orderedItemApp := NewOrderedItemApplication(a.p)
	orderedItemApp.SaveRawOrderItems(rawOrder.Products, int64(order.ID))
	
	if err != nil {
		// Handle error if needed
		return nil, map[string]string{"error": "failed to save order"}
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


