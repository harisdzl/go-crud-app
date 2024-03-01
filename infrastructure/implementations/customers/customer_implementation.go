package customers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/customer_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	p *base.Persistence
	c *context.Context
}

func NewCustomerRepository(p *base.Persistence, c *context.Context) *CustomerRepo {
	return &CustomerRepo{p, c}
}

var _ customer_repository.CustomerRepository = &CustomerRepo{}

func (c *CustomerRepo) SaveCustomer(customer *customer_entity.Customer) (*customer_entity.Customer, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", c.p)

	dbErr := map[string]string{}
	err := c.p.DB.Debug().Create(&customer).Error
	if err != nil {
		fmt.Println("Failed to create customer")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	cacheRepo.SetKey(fmt.Sprintf("%v_CUSTOMER", customer.ID), customer, time.Minute * 15)
	
	return customer, nil
}


func (c *CustomerRepo) GetCustomer(id int64) (*customer_entity.Customer, error) {
	var customer *customer_entity.Customer

	cacheRepo := cache.NewCacheRepository("Redis", c.p)
	_ = cacheRepo.GetKey(fmt.Sprintf("%v_CUSTOMER", id), &customer)
	if customer == nil {
		err := c.p.DB.Debug().Where("id = ?", id).Take(&customer).Error
		if err != nil {
			fmt.Println("Failed to get customer")
		}
		if customer != nil && customer.ID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprintf("%v_CUSTOMER", id), customer, time.Minute * 15)
		}
	}


	return customer, nil
}

func (c *CustomerRepo) GetAllCustomers() ([]customer_entity.Customer, error) {
	var customers []customer_entity.Customer
	err := c.p.DB.Debug().Find(&customers).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return customers, nil
}

func (c *CustomerRepo) UpdateCustomer(customer *customer_entity.Customer) (*customer_entity.Customer, error) {
	cacheRepo := cache.NewCacheRepository("Redis", c.p)


	err := c.p.DB.Debug().Where("id = ?", customer.ID).Updates(&customer).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = cacheRepo.SetKey(fmt.Sprintf("%v_CUSTOMER", customer.ID), customer, time.Minute * 15)

	return customer, nil
}

func (c *CustomerRepo) DeleteCustomer(id int64) error {
	var customer customer_entity.Customer

	err := c.p.DB.Debug().Where("id = ?", id).Delete(&customer).Error
	
	cacheRepo := cache.NewCacheRepository("Redis", c.p)

	cacheRepo.DelKey(fmt.Sprintf("%v_CUSTOMER", id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}