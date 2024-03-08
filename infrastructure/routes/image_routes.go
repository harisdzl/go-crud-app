package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func ImageRoutes(router *gin.RouterGroup, p *base.Persistence) {
    images := handlers.NewImage(p)
       
    router.POST("admin/images", images.SaveImage)
    router.GET("admin/images/:image_id", images.GetImage)
    router.GET("admin/images/products/:product_id", images.GetAllImagesOfProduct)
    router.DELETE("admin/images/:image_id", images.DeleteImage)
}
