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
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type Order struct {
	OrderRepo   order_repository.OrderHandlerRepository
	Persistence *base.Persistence
}

func NewOrder(p *base.Persistence) *Order {
	return &Order{
		Persistence: p,
	}
}

func (or Order) SaveOrder(c *gin.Context) {
	// Start a new span for the handler function
	channels := []string{"Zap", "Honeycomb"}
	loggerRepo, loggerErr := logger.NewLoggerRepository(channels, or.Persistence, c, "handlers/SaveOrder")
	if loggerErr != nil {
		loggerRepo.Warn("Error in initializing logger", map[string]interface{}{})
	}
	loggerRepo.SetContextWithSpan()
	defer loggerRepo.End()


	
	responseContextData := entity.ResponseContext{Ctx: c}
	rawOrder := order_entity.RawOrder{}
	userIdString := c.GetString("userID")
	userId, userIdErr := strconv.ParseInt(userIdString, 10, 64)

	if userIdErr != nil {
		// Log error within the span
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid user", ""))
		return
	}

	rawOrder.CustomerID = userId
	if err := c.ShouldBindJSON(&rawOrder); err != nil {
		// Log error within the span
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}


	or.OrderRepo = application.NewOrderApplication(or.Persistence, c)
	savedOrder, saveErr := or.OrderRepo.SaveOrderFromRaw(rawOrder)
	if saveErr != nil {
		// Log error within the span
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, saveErr.Error(), ""))
		return
	}

	results := map[string]interface{}{
		"user_ip" : c.ClientIP(),
		"user_id" : userId,
		"json_data": savedOrder,
	}

	loggerRepo.Info("Order saved successfully", results)
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Order saved successfully", &savedOrder))
}

//	@Summary		Get All Orders
//	@Description	Retrieves all orders.
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ResponseContext	"Success"
//	@Failure		500	{object}	entity.ResponseContext	"Internal server error"
//	@Router			/orders [get]
func (or Order) GetAllOrders(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	

	or.OrderRepo = application.NewOrderApplication(or.Persistence, c)

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

// GetOrder retrieves a specific order by its ID.
//	@Summary		Get Order
//	@Description	Retrieves a specific order by its ID.
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			order_id	path		int						true	"Order ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/orders/{order_id} [get]
func (or Order) GetOrder(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	or.OrderRepo = application.NewOrderApplication(or.Persistence, c)

	order, err := or.OrderRepo.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Order %v obtained", orderID), order))
}

// DeleteOrder deletes a specific order by its ID.
//	@Summary		Delete Order
//	@Description	Deletes a specific order by its ID.
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			order_id	path		int						true	"Order ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/orders/{order_id} [delete]
func (or Order) DeleteOrder(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	or.OrderRepo = application.NewOrderApplication(or.Persistence, c)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	deleteErr := or.OrderRepo.DeleteOrder(orderID)
	// TODO: when deleting a order, need to delete all the inventory in it
	orderedItems, orderedItemsErr := application.NewOrderedItemApplication(or.Persistence, c).GetAllOrderedItemsForOrder(orderID)
	if orderedItemsErr != nil {	
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	application.NewOrderedItemApplication(or.Persistence, c).ReverseOrder(orderedItems)


	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Order deleted successfully", ""))
}

// UpdateOrder updates a specific order by its ID.
//	@Summary		Update Order
//	@Description	Updates a specific order by its ID.
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			order_id	path		int						true	"Order ID"
//	@Param			Order		body		order_entity.Order			true	"Order object to be updated"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		404			{object}	entity.ResponseContext	"Order not found"
//	@Failure		422			{object}	entity.ResponseContext	"Unprocessable entity"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/orders/{order_id} [put]
func (or Order) UpdateOrder(c *gin.Context) {
	// Start a new span for the handler function
	channels := []string{"Zap", "Honeycomb"}
	loggerRepo, loggerErr := logger.NewLoggerRepository(channels, or.Persistence, c, "handlers/UpdateOrder")
	loggerRepo.SetContextWithSpan()

	defer loggerRepo.End()
	if loggerErr != nil {
		loggerRepo.Warn("Error in initializing logger", map[string]interface{}{})
	}
	responseContextData := entity.ResponseContext{Ctx: c}
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid Order ID", ""))
		return
	}

	// Check if the Order exists
	or.OrderRepo = application.NewOrderApplication(or.Persistence, c)

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
	
	loggerRepo.SetContextWithSpan()
	or.OrderRepo = application.NewOrderApplication(or.Persistence, c)

	// Update the Order
	updatedOrder, updateErr := or.OrderRepo.UpdateOrder(existingOrder)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Order updated successfully", updatedOrder))
}
