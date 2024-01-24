package product_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
)

type ProductRepository interface {
	SaveProduct(*entity.Product) (*entity.Product, map[string]string)
	SaveMultipleProducts(*[]entity.Product) (*[]entity.Product, map[string]string)
	GetProduct(int64) (*entity.Product, error)
	GetAllProducts() ([]entity.Product, error)
	UpdateProduct(*entity.Product) (*entity.Product, error)
	DeleteProduct(int64) error
	SearchProduct(string) ([]entity.Product, error)
	UpdateProductsInSearchDB() (error)
}