package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func ImageRoutes(router *gin.Engine, p *base.Persistence) {
    images := handlers.NewImage(p)
       
    router.POST("/images", images.SaveImage)
    router.GET("/images/:image_id", images.GetImage)
    router.GET("/images/products/:product_id", images.GetAllImagesOfProduct)
    router.DELETE("/images/:image_id", images.DeleteImage)
}
