package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func WarehouseRoutes(router *gin.RouterGroup, p *base.Persistence) {
    warehouses := handlers.NewWarehouse(p)

    
    router.POST("admin/warehouses", warehouses.SaveWarehouse)
    // router.POST("admin/warehouses/multiple", warehouses.SaveMultiplewarehouses)
    router.GET("admin/warehouses", warehouses.GetAllWarehouses)
    router.GET("admin/warehouses/:warehouse_id", warehouses.GetWarehouse)
    router.GET("admin/warehouses/:warehouse_id/inventories", warehouses.GetInventoriesInWarehouse)
    router.PUT("admin/warehouses/:warehouse_id", warehouses.UpdateWarehouse)
    router.DELETE("admin/warehouses/:warehouse_id", warehouses.DeleteWarehouse)
    router.GET("admin/warehouses/search", warehouses.SearchWarehouse)
    router.POST("admin/warehouses/search", warehouses.UpdateWarehouseSearchDB)
}
