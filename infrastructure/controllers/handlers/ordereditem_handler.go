package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/ordereditem_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type OrderedItem struct {
	OrderedItemRepo ordereditem_repository.OrderedItemRepository
	Persistence *base.Persistence
}



func NewOrderedItem(p *base.Persistence) *OrderedItem {
	return &OrderedItem{
		Persistence: p,
	}
}


// GetAllOrderedItems retrieves all ordered items.
//	@Summary		Get All Ordered Items
//	@Description	Retrieves all ordered items.
//	@Tags			OrderedItem
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ResponseContext	"Success"
//	@Failure		500	{object}	entity.ResponseContext	"Internal server error"
//	@Router			/ordereditems [get]
func (or OrderedItem) GetAllOrderedItems(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	or.OrderedItemRepo = application.NewOrderedItemApplication(or.Persistence)

	allOrderedItems, err := or.OrderedItemRepo.GetAllOrderedItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	results := map[string]interface{}{
		"results" : allOrderedItems,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "All orders obtained successfully", results))
}

// GetAllOrderedItemsForOrder retrieves all ordered items for a specific order.
//	@Summary		Get All Ordered Items For Order
//	@Description	Retrieves all ordered items for a specific order.
//	@Tags			OrderedItem
//	@Accept			json
//	@Produce		json
//	@Param			order_id	path		int						true	"Order ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/ordereditems/{order_id} [get]
func (or OrderedItem) GetAllOrderedItemsForOrder(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	orderId, err := strconv.ParseInt(c.Param("order_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	or.OrderedItemRepo = application.NewOrderedItemApplication(or.Persistence)

	orderedItems, err := or.OrderedItemRepo.GetAllOrderedItemsForOrder(orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Ordered Items for order %v obtained", orderId), orderedItems))
}