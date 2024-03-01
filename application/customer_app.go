package application

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/customer_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/customers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type customerApp struct {
	p *base.Persistence
	c *gin.Context
}

func NewCustomerApplication(p *base.Persistence, c *gin.Context) customer_repository.CustomerRepository {
	return &customerApp{p, c}
}

func (a *customerApp) SaveCustomer(customer *customer_entity.Customer) (*customer_entity.Customer, map[string]string) {
	repocustomer := customers.NewCustomerRepository(a.p, a.c)
	return repocustomer.SaveCustomer(customer)
}

func (a *customerApp) GetCustomer(customerId int64) (*customer_entity.Customer, error) {
	repocustomer := customers.NewCustomerRepository(a.p, a.c)
	return repocustomer.GetCustomer(customerId)
}

func (a *customerApp) GetAllCustomers() ([]customer_entity.Customer, error) {
	repocustomer := customers.NewCustomerRepository(a.p, a.c)
	return repocustomer.GetAllCustomers()
}
	
func (a *customerApp) UpdateCustomer(customer *customer_entity.Customer) (*customer_entity.Customer, error) {
	repocustomer := customers.NewCustomerRepository(a.p, a.c)
	return repocustomer.UpdateCustomer(customer)
}

func (a *customerApp) DeleteCustomer(customerId int64) error {
	repocustomer := customers.NewCustomerRepository(a.p, a.c)
	return repocustomer.DeleteCustomer(customerId)
}


