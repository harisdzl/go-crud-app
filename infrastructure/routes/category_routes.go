package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func CategoryRoutes(router *gin.RouterGroup, p *base.Persistence) {
    categories := handlers.NewCategory(p)
       
    router.POST("admin/categories", categories.SaveCategory)
    router.GET("admin/categories/:category_id", categories.GetCategory)
    router.GET("admin/categories/parents/:category_id", categories.GetParentCategories)
    router.GET("admin/categories", categories.GetAllCategories)
	router.PUT("admin/categories/:category_id", categories.UpdateCategory)
    router.DELETE("admin/categories/:category_id", categories.DeleteCategory)
}
