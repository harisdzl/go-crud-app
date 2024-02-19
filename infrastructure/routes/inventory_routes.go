package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func InventoryRoutes(router *gin.Engine, p *base.Persistence) {
    inventories := handlers.NewInventory(p)
       
    router.GET("/products/:product_id/inventories", inventories.GetInventory)
    router.PUT("/products/:product_id/inventories", inventories.UpdateInventory)
}