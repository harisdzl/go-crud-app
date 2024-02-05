package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func OrderRoutes(router *gin.Engine, p *base.Persistence) {
    orders := handlers.NewOrder(p)
    
    router.POST("/orders", orders.SaveOrder)
    router.GET("/orders", orders.GetAllOrders)
    router.GET("/orders/:order_id", orders.GetOrder)
    router.PUT("/orders/:order_id", orders.UpdateOrder)
    router.DELETE("/orders/:order_id", orders.DeleteOrder)
}
