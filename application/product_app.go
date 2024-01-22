package application

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
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
	SaveProduct(*entity.Product) (*entity.Product, map[string]string)
	SaveMultipleProducts(*[]entity.Product) (*[]entity.Product, map[string]string)
	GetProduct(int64) (*entity.Product, error)
	GetAllProducts() ([]entity.Product, error)
	UpdateProduct(*entity.Product) (*entity.Product, error)
	DeleteProduct(int64) error
	SearchProduct(string) ([]interface{}, error)
	UpdateProductsInMongo() (error)
}

func (a *productApp) SaveProduct(product *entity.Product) (*entity.Product, map[string]string) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.SaveProduct(product)
}

func (a *productApp) SaveMultipleProducts(productList *[]entity.Product) (*[]entity.Product, map[string]string) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.SaveMultipleProducts(productList)
}

func (a *productApp) GetProduct(productId int64) (*entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.GetProduct(productId)
}

func (a *productApp) GetAllProducts() ([]entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.GetAllProducts()
}
	
func (a *productApp) UpdateProduct(product *entity.Product) (*entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.UpdateProduct(product)
}

func (a *productApp) DeleteProduct(productId int64) error {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.DeleteProduct(productId)
}

func (a *productApp) SearchProduct(name string) ([]entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.SearchProduct(name)
}

func (a *productApp) UpdateProductsInMongo() (error) {
	repoProduct := products.NewProductRepository(a.p)
	return repoProduct.UpdateProductsInMongo()
}



