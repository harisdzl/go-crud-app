package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func WarehouseRoutes(router *gin.RouterGroup, p *base.Persistence) {
    warehouses := handlers.NewWarehouse(p)

    
    router.POST("/warehouses", warehouses.SaveWarehouse)
    // router.POST("/warehouses/multiple", warehouses.SaveMultiplewarehouses)
    router.GET("/warehouses", warehouses.GetAllWarehouses)
    router.GET("/warehouses/:warehouse_id", warehouses.GetWarehouse)
    router.GET("/warehouses/:warehouse_id/inventories", warehouses.GetInventoriesInWarehouse)
    router.PUT("/warehouses/:warehouse_id", warehouses.UpdateWarehouse)
    router.DELETE("/warehouses/:warehouse_id", warehouses.DeleteWarehouse)
    router.GET("/warehouses/search", warehouses.SearchWarehouse)
    router.POST("/warehouses/search", warehouses.UpdateWarehouseSearchDB)
}
