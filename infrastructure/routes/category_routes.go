package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func CategoryRoutes(router *gin.Engine, p *base.Persistence) {
    categories := handlers.NewCategory(p)
       
    router.POST("/categories", categories.SaveCategory)
    router.GET("/categories/:category_id", categories.GetCategory)
    router.GET("/categories/parents/:category_id", categories.GetParentCategories)
    router.GET("/categories", categories.GetAllCategories)
	router.PUT("/categories/:category_id", categories.UpdateCategory)
    router.DELETE("/categories/:category_id", categories.DeleteCategory)
}
