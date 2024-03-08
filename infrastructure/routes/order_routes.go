package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func OrderRoutes(router *gin.RouterGroup, p *base.Persistence) {
    orders := handlers.NewOrder(p)
    
    router.POST("admin/orders", orders.SaveOrder)
    router.GET("admin/orders", orders.GetAllOrders)
    router.GET("admin/orders/:order_id", orders.GetOrder)
    router.PUT("admin/orders/:order_id", orders.UpdateOrder)
    router.DELETE("admin/orders/:order_id", orders.DeleteOrder)
}
