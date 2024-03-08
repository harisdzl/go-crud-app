package application

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/product_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/products"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search"
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
	channels := []string{"Zap", "Honeycomb"}
	loggerRepo, loggerErr := logger.NewLoggerRepository(channels, a.p, a.c, "application/SearchProduct")
	if loggerErr != nil {
		return nil, loggerErr
	}
	loggerRepo.SetContextWithSpan()
	defer loggerRepo.End()

	searchProvider := os.Getenv("SEARCH_PROVIDER")

	repoSearch := search.NewSearchRepository(searchProvider, a.p, a.c)

	repoProduct := products.NewProductRepository(a.p, a.c)

	indexName := "products"
	// Extract the results from the cursor
	var results []map[string]interface{}
	var searchProducts []product_entity.Product
	err := repoSearch.SearchDocByName(name, indexName, &results)

	for _, result := range results {
		productId, productIdErr := strconv.ParseInt(result["id"].(string), 10, 64)

		if productIdErr != nil {
			loggerRepo.Error("Error in converting product id", map[string]interface{}{"data": result})
			return nil, productIdErr
		}

		product, productErr := repoProduct.GetProduct(productId)
		if productErr != nil {
			return nil, productErr
		}

		searchProducts = append(searchProducts, *product)
	}
	if err != nil {
		fmt.Println(err)
	}

	if len(results) == 0 {
		fmt.Println("No such product of name: " + name)
	}
	log.Println(results)
	return searchProducts, nil
}

func (a *productApp) UpdateProductsInSearchDB() (error) {
	searchProvider := os.Getenv("SEARCH_PROVIDER")
	searchRepo := search.NewSearchRepository(searchProvider, a.p, a.c)
	collectionName := "products"

	products, err := a.GetAllProducts()
	
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var allProducts []interface{}

    for _, p := range products {
		idString := fmt.Sprint(p.ID)
		searchProduct := map[string]interface{}{
			"id" : idString,
			"name" : p.Name,
			"description" : p.Description,
			"category" : p.Category.Name,
		}

        allProducts = append(allProducts, searchProduct)
    }


	searchDeleteAllErr := searchRepo.DeleteAllDoc(collectionName, allProducts)
	searchInsertAll := searchRepo.InsertAllDoc(collectionName, allProducts)

	if searchDeleteAllErr != nil {
		return errors.New("Fail to delete all docs")
	}

	if searchInsertAll != nil {
		return errors.New("Fail to update search db with all products")
	}

	return nil
}



