package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/order_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/order_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type Order struct {
	OrderRepo order_repository.OrderHandlerRepository
	Persistence *base.Persistence
}



func NewOrder(p *base.Persistence) *Order {
	return &Order{
		Persistence: p,
	}
}

// SaveOrder saves a single Order to the database.
// @Summary Save a single Order
// @Description SaveOrder saves a single Order to the database.
// @Tags Order
// @Accept json
// @Produce json
// @Param Order body entity.Order true "Order object to be saved"
// @Success 201 {object} entity.Order "Successfully saved Order"
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 422 {object} map[string]string "Unprocessable entity"
// @Router /Orders [post]
func (or Order) SaveOrder(c *gin.Context) {
    responseContextData := entity.ResponseContext{Ctx: c}
    rawOrder := order_entity.RawOrder{}
    if err := c.ShouldBindJSON(&rawOrder); err != nil {
		log.Println(err)
        c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
        return
    }

	or.OrderRepo = application.NewOrderApplication(or.Persistence)

    savedOrder, saveErr := or.OrderRepo.SaveOrderFromRaw(rawOrder)
    if saveErr != nil {
        c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, saveErr["db_error"], ""))
        return
    }

    c.JSON(http.StatusCreated, responseContextData.ResponseData(entity.StatusSuccess, "Order saved successfully", savedOrder))
}

func (or Order) GetAllOrders(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	or.OrderRepo = application.NewOrderApplication(or.Persistence)

	allOrders, err := or.OrderRepo.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	results := map[string]interface{}{
		"results" : allOrders,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "All orders obtained successfully", results))
}

func (or Order) GetOrder(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	or.OrderRepo = application.NewOrderApplication(or.Persistence)

	order, err := or.OrderRepo.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Order %v obtained", orderID), order))
}

func (or Order) DeleteOrder(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	or.OrderRepo = application.NewOrderApplication(or.Persistence)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	deleteErr := or.OrderRepo.DeleteOrder(orderID)
	// TODO: when deleting a order, need to delete all the inventory in it
	orderedItems, orderedItemsErr := application.NewOrderedItemApplication(or.Persistence).GetAllOrderedItemsForOrder(orderID)
	if orderedItemsErr != nil {	
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	application.NewOrderedItemApplication(or.Persistence).ReverseOrder(orderedItems)


	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Order deleted successfully", ""))
}

func (or Order) UpdateOrder(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid Order ID", ""))
		return
	}

	// Check if the Order exists
	or.OrderRepo = application.NewOrderApplication(or.Persistence)

	existingOrder, err := or.OrderRepo.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, responseContextData.ResponseData(entity.StatusFail, "Order not found", ""))
		return
	}

	// Bind the JSON request body to the existing Order
	if err := c.ShouldBindJSON(&existingOrder); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	or.OrderRepo = application.NewOrderApplication(or.Persistence)

	// Update the Order
	updatedOrder, updateErr := or.OrderRepo.UpdateOrder(existingOrder)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Order updated successfully", updatedOrder))
}
