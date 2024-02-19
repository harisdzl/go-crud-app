package product_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
)

type ProductRepository interface {
	SaveProduct(*product_entity.Product) (*product_entity.Product, map[string]string)
	// SaveMultipleProducts(*[]product_entity.Product) (*[]product_entity.Product, map[string]string)
	GetProduct(int64) (*product_entity.Product, error)
	GetAllProducts() ([]product_entity.Product, error)
	UpdateProduct(*product_entity.Product) (*product_entity.Product, error)
	DeleteProduct(int64) error
	SearchProduct(string) ([]product_entity.Product, error)
	UpdateProductsInSearchDB() (error)
}


type ProductHandlerRepository interface {
	SaveProductAndInventory(product_entity.ProductForInventory) (*product_entity.Product, *inventory_entity.Inventory, map[string]string)
	// SaveMultipleProducts(*[]product_entity.Product) (*[]product_entity.Product, map[string]string)
	GetProduct(int64) (*product_entity.Product, error)
	GetAllProducts() ([]product_entity.Product, error)
	UpdateProduct(*product_entity.Product) (*product_entity.Product, error)
	DeleteProduct(int64) error
	SearchProduct(string) ([]product_entity.Product, error)
	UpdateProductsInSearchDB() (error)
}