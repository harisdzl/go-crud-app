package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func ProductRoutes(router *gin.Engine, p *base.Persistence) {
    products := handlers.NewProduct(p)

    products.UpdateMongo()

    router.POST("/products", products.SaveProduct)
    router.POST("/products/multiple", products.SaveMultipleProducts)
    router.GET("/products", products.GetAllProducts)
    router.GET("/products/:product_id", products.GetProduct)
    router.PUT("/products/:product_id", products.UpdateProduct)
    router.DELETE("/products/:product_id", products.DeleteProduct)
    router.GET("/products/search", products.SearchProduct)
}
