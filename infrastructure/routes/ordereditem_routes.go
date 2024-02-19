package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func OrderedItemRoutes(router *gin.Engine, p *base.Persistence) {
    orderedItems := handlers.NewOrderedItem(p)
    
    router.GET("/ordereditems", orderedItems.GetAllOrderedItems)
    router.GET("/orders/:order_id/ordereditems", orderedItems.GetAllOrderedItemsForOrder)
}
