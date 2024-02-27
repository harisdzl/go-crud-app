package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/docs"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/middleware"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitMiddleware(engine *gin.Engine) {
	engine.Use(middleware.CORSMiddleware())
}

func InitRouter(p *base.Persistence) *gin.Engine {
    r := gin.Default()
    p.Automigrate()
    InitMiddleware(r)
    AuthRoutesPublic(r.Group("/"), p)
    CustomerPublicRoutes(r.Group("/"), p)
    private := r.Group("/")
    authMiddleware := middleware.AuthHandler(p)


    // Apply Honeycomb middleware to all routes
    honeycombMiddleware := middleware.HoneycombHandler()


    private.Use(authMiddleware, honeycombMiddleware)

    // Define routes within the private group
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

    // Swagger documentation setup
    docs.SwaggerInfo.Title = "Quqo Challenge"
    docs.SwaggerInfo.BasePath = "/"
    docs.SwaggerInfo.Description = "Documentation - Quqo Challenge"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Schemes = []string{"http", "https"}

    // Use ginSwagger middleware to serve the API docs
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return r
}
