package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func AuthRoutesPublic(router *gin.RouterGroup, p *base.Persistence) {
    auth := handlers.NewAuth(p)
	
	router.POST("admin/login", auth.Login)
}

func AuthRoutesPrivate(router *gin.RouterGroup, p *base.Persistence) {
    auth := handlers.NewAuth(p)
	
	router.POST("admin/logout", auth.Logout)
}
