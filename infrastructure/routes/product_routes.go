package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func ProductRoutes(router *gin.RouterGroup, p *base.Persistence) {
    products := handlers.NewProduct(p)

    
    router.POST("admin/products", products.SaveProduct)
    // router.POST("admin/products/multiple", products.SaveMultipleProducts)
    router.GET("admin/products", products.GetAllProducts)
    router.GET("admin/products/:product_id", products.GetProduct)
    router.PUT("admin/products/:product_id", products.UpdateProduct)
    router.DELETE("admin/products/:product_id", products.DeleteProduct)
    router.GET("admin/products/search", products.SearchProduct)
    router.POST("admin/products/search", products.UpdateProductSearchDB)
}
