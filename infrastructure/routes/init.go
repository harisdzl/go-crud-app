package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/docs"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/middleware"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func InitRouter(p *base.Persistence) *gin.Engine {

	r := gin.Default()
	p.Automigrate()
	AuthRoutesPublic(r.Group("/"), p)
	CustomerPublicRoutes(r.Group("/"), p)
	private := r.Group("/")
	middleware := middleware.AuthHandler(p)
	private.Use(middleware)

	{
		ProductRoutes(private, p)
		InventoryRoutes(private, p)
		WarehouseRoutes(private, p)
		ImageRoutes(private, p)
		CategoryRoutes(private, p)
		CustomerPrivateRoutes(private, p)
		OrderRoutes(private, p)
		OrderedItemRoutes(private, p)
		AuthRoutesPrivate(private, p)
	}

	docs.SwaggerInfo.Title = "Quqo Challenge"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Description = "Documentation - Quqo Challenge"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	
	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	return r
}