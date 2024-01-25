package application

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/product_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/products"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type productApp struct {
	p *base.Persistence
}

func NewProductApplication(p *base.Persistence) product_repository.ProductRepository {
	return &productApp{p}
}

type ProductAppInterface interface {
	SaveProduct(*product_entity.Product) (*product_entity.Product, map[string]string)
	// SaveMultipleProducts(*[]product_entity.Product) (*[]product_entity.Product, map[string]string)
	GetProduct(int64) (*product_entity.Product, error)
	GetAllProducts() ([]product_entity.Product, error)
	UpdateProduct(*product_entity.Product) (*product_entity.Product, error)
	DeleteProduct(int64) error
	SearchProduct(string) ([]interface{}, error)
	UpdateProductsInSearchDB() (error)
}

func (a *productApp) SaveProduct(product *product_entity.Product) (*product_entity.Product, map[string]string) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.SaveProduct(product)
}

// func (a *productApp) SaveMultipleProducts(productList *[]product_entity.Product) (*[]product_entity.Product, map[string]string) {
// 	repoProduct := products.NewProductRepository(a.p)
// 	return repoProduct.SaveMultipleProducts(productList)
// }

func (a *productApp) GetProduct(productId int64) (*product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.GetProduct(productId)
}

func (a *productApp) GetAllProducts() ([]product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.GetAllProducts()
}
	
func (a *productApp) UpdateProduct(product *product_entity.Product) (*product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.UpdateProduct(product)
}

func (a *productApp) DeleteProduct(productId int64) error {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.DeleteProduct(productId)
}

func (a *productApp) SearchProduct(name string) ([]product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.SearchProduct(name)
}

func (a *productApp) UpdateProductsInSearchDB() (error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.UpdateProductsInSearchDB()
}



