package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/controllers/handlers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

func CustomerRoutes(router *gin.Engine, p *base.Persistence) {
    customers := handlers.NewCustomer(p)
    
    router.POST("/customers", customers.SaveCustomer)
    router.GET("/customers", customers.GetAllCustomers)
    router.GET("/customers/:customer_id", customers.GetCustomer)
    router.PUT("/customers/:customer_id", customers.UpdateCustomer)
    router.DELETE("/customers/:customer_id", customers.DeleteCustomer)
}