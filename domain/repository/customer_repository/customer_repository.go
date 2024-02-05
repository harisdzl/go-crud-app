package customer_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
)

type CustomerRepository interface {
	SaveCustomer(*customer_entity.Customer) (*customer_entity.Customer, map[string]string)
	GetCustomer(int64) (*customer_entity.Customer, error)
	GetAllCustomers() ([]customer_entity.Customer, error)
	UpdateCustomer(*customer_entity.Customer) (*customer_entity.Customer, error)
	DeleteCustomer(int64) error
}