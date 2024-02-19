package application

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/customer_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/customers"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type customerApp struct {
	p *base.Persistence
}

func NewCustomerApplication(p *base.Persistence) customer_repository.CustomerRepository {
	return &customerApp{p}
}

func (a *customerApp) SaveCustomer(customer *customer_entity.Customer) (*customer_entity.Customer, map[string]string) {
	repocustomer := customers.NewCustomerRepository(a.p)
	return repocustomer.SaveCustomer(customer)
}

func (a *customerApp) GetCustomer(customerId int64) (*customer_entity.Customer, error) {
	repocustomer := customers.NewCustomerRepository(a.p)
	return repocustomer.GetCustomer(customerId)
}

func (a *customerApp) GetAllCustomers() ([]customer_entity.Customer, error) {
	repocustomer := customers.NewCustomerRepository(a.p)
	return repocustomer.GetAllCustomers()
}
	
func (a *customerApp) UpdateCustomer(customer *customer_entity.Customer) (*customer_entity.Customer, error) {
	repocustomer := customers.NewCustomerRepository(a.p)
	return repocustomer.UpdateCustomer(customer)
}

func (a *customerApp) DeleteCustomer(customerId int64) error {
	repocustomer := customers.NewCustomerRepository(a.p)
	return repocustomer.DeleteCustomer(customerId)
}


