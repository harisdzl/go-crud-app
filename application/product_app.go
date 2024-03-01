package application

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/product_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/products"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type productApp struct {
	p *base.Persistence
	c *gin.Context
}

func NewProductApplication(p *base.Persistence, c *gin.Context) product_repository.ProductHandlerRepository {
	return &productApp{p, c}
}

func ConvertProductandInventory(productForInventory product_entity.ProductForInventory) (product_entity.Product, inventory_entity.Inventory){
	var product product_entity.Product
	var inventory inventory_entity.Inventory

	product.ID = productForInventory.ID
	product.Name = productForInventory.Name
	product.Description = productForInventory.Description
	product.Price = productForInventory.Price
	product.CategoryID = productForInventory.CategoryID

	inventory.ProductID = product.ID
	inventory.WarehouseID = productForInventory.WarehouseID
	inventory.Stock = productForInventory.Stock

	return product, inventory

}
func (a *productApp) SaveProductAndInventory(productForInventory product_entity.ProductForInventory) (*product_entity.Product, *inventory_entity.Inventory, map[string]string) {
    product, inventory := ConvertProductandInventory(productForInventory)
    repoProduct := products.NewProductRepository(a.p, a.c)
    repoInventory := inventories.NewInventoryRepository(a.p, a.c)
    savedProduct, saveErr := repoProduct.SaveProduct(&product)
    if saveErr != nil {
        return nil, nil, saveErr
    }
    inventory.ProductID = savedProduct.ID // Set ProductID to the ID of the newly created product
    _, saveInventoryErr := repoInventory.SaveInventory(&inventory)
    if saveInventoryErr != nil {
        return nil, nil, saveInventoryErr
    }
    return savedProduct, &inventory, nil
}


func (a *productApp) GetProduct(productId int64) (*product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p, a.c)
	return repoProduct.GetProduct(productId)
}

func (a *productApp) GetAllProducts() ([]product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p, a.c)
	return repoProduct.GetAllProducts()
}
	
func (a *productApp) UpdateProduct(product *product_entity.Product) (*product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p, a.c)
	return repoProduct.UpdateProduct(product)
}

func (a *productApp) DeleteProduct(productId int64) error {
	repoProduct := products.NewProductRepository(a.p, a.c)
	return repoProduct.DeleteProduct(productId)
}

func (a *productApp) SearchProduct(name string) ([]product_entity.Product, error) {
	repoProduct := products.NewProductRepository(a.p, a.c)
	return repoProduct.SearchProduct(name)
}

func (a *productApp) UpdateProductsInSearchDB() (error) {
	repoProduct := products.NewProductRepository(a.p, a.c)
	return repoProduct.UpdateProductsInSearchDB()
}



