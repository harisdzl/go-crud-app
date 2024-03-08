package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func CustomerPrivateRoutes(router *gin.RouterGroup, p *base.Persistence) {
    customers := handlers.NewCustomer(p)
    
    router.GET("admin/customers", customers.GetAllCustomers)
    router.GET("admin/customers/:customer_id", customers.GetCustomer)
    router.PUT("admin/customers/:customer_id", customers.UpdateCustomer)
    router.DELETE("admin/customers/:customer_id", customers.DeleteCustomer)
}


func CustomerPublicRoutes(router *gin.RouterGroup, p *base.Persistence) {
    customers := handlers.NewCustomer(p)
    
    router.POST("admin/customers", customers.SaveCustomer)
}
