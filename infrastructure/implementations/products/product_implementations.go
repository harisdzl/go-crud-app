package products

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/product_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/logger"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

// To manage new product repositories in the database

// Product Repository struct
type ProductRepo struct {
	p *base.Persistence
	c *gin.Context
}

func NewProductRepository(p *base.Persistence, c *gin.Context) *ProductRepo {
	return &ProductRepo{p, c}
}

// To explicitly check that the ProductRepo implements the repository.ProductRepository interface
var _ product_repository.ProductRepository = &ProductRepo{}

func (r *ProductRepo) SaveProduct(product *product_entity.Product) (*product_entity.Product, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	searchRepo := search.NewSearchRepository("Mongo", r.p, r.c)
	
	dbErr := map[string]string{}
	err := r.p.DB.Debug().Create(&product).Error
	collectionName := "products"
	if err != nil {
		fmt.Println("Failed to create product")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	fmt.Printf("Type of product: %T\n", product) // Log the type of the product


	searchErr := searchRepo.InsertDoc(collectionName, interface{}(product))

	if searchErr != nil {
		fmt.Println(searchErr)
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	cacheRepo.SetKey(fmt.Sprintf("%v_PRODUCTS", product.ID), product, time.Minute * 15)
	
	return product, nil
}

// func (r *ProductRepo) SaveMultipleProducts(product *[]product_entity.Product) (*[]product_entity.Product, map[string]string) {
// 	dbErr := map[string]string{}
// 	err := r.p.DB.Debug().Create(&product).Error
// 	if err != nil {
// 		fmt.Println("Failed to create products")
// 		fmt.Println(err)
// 		dbErr["db_error"] = "database error"
// 		return nil, dbErr
// 	}
// 	collectionName := "products"

// 	var allProducts []interface{}

//     for _, p := range *product {
//         allProducts = append(allProducts, p)
//     }

//     searchRepo := search.NewSearchRepository("Mongo", r.p)

//     searchErr := searchRepo.InsertAllDoc(collectionName, allProducts)

//     if searchErr != nil {
//         fmt.Println("Failed to insert products into search index")
//         fmt.Println(searchErr)
//         dbErr["search_error"] = "search index error"
//         return nil, dbErr
//     }



// 	return product, nil
// }

func (r *ProductRepo) GetProduct(id int64) (*product_entity.Product, error) {
    channels := []string{"Zap", "Honeycomb"}
    loggerRepo, loggerErr := logger.NewLoggerRepository(channels, r.p, r.c, "implementations/GetProduct")

    if loggerErr != nil {
        return nil, loggerErr
    }
    defer loggerRepo.Span.End()
    
    var product *product_entity.Product

    cacheRepo := cache.NewCacheRepository("Redis", r.p)
    _ = cacheRepo.GetKey(fmt.Sprintf("%v_PRODUCTS", id), &product)
    if product == nil {
        err := r.p.DB.Debug().
		Preload("Category").
		Preload("Images").
		Preload("Inventory").
		Where("id = ?", id).Take(&product).Error
        if err != nil {
            fmt.Println("Failed to get product")
            return nil, err
        }
        if product != nil && product.ID > 0 {
            _ = cacheRepo.SetKey(fmt.Sprintf("%v_PRODUCTS", id), product, time.Minute * 15)
        }
    }

    loggerRepo.Info("Product retrieved", map[string]interface{}{"data": product})

    return product, nil
}


func (r *ProductRepo) GetAllProducts() ([]product_entity.Product, error) {
	var products []product_entity.Product
	err := r.p.DB.Debug().
	Preload("Category").
	Preload("Images").
	Preload("Inventory").
	Find(&products).Error

	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepo) UpdateProduct(product *product_entity.Product) (*product_entity.Product, error) {
	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	// searchRepo := search.NewSearchRepository("Mongo", r.p)
	// collectionName := "products"

	err := r.p.DB.Debug().Where("id = ?", product.ID).Updates(&product).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	_ = cacheRepo.SetKey(fmt.Sprintf("%v_PRODUCTS", product.ID), product, time.Minute * 15)

	// searchErr := searchRepo.UpdateDoc(uint(product.ID), collectionName, &product)

	// if searchErr != nil {
	// 	return nil, err
	// }

	// if errors.Is(searchErr, gorm.ErrRecordNotFound) {
	// 	return nil, err
	// }



	return product, nil
}

func (r *ProductRepo) DeleteProduct(id int64) error {
	var product product_entity.Product
	searchRepo := search.NewSearchRepository("Mongo", r.p, r.c)
	collectionName := "products"
	fieldName := "id"
	err := r.p.DB.Debug().Where("id = ?", id).Delete(&product).Error
	
	searchErr := searchRepo.DeleteSingleDoc(fieldName, collectionName, id)
	cacheRepo := cache.NewCacheRepository("Redis", r.p)

	cacheRepo.DelKey(fmt.Sprintf("%v_PRODUCTS", id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	if searchErr != nil {
		return errors.New("database error, please try again")
	}

	return nil
}



func (r *ProductRepo) UpdateProductsInSearchDB() error {
	searchRepo := search.NewSearchRepository("Mongo", r.p, r.c)
	collectionName := "products"

	products, err := r.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var allProducts []interface{}

	for _, p := range products {
		allProducts = append(allProducts, p)
	}
	searchDeleteAllErr := searchRepo.DeleteAllDoc(collectionName, allProducts)
	searchInsertAll := searchRepo.InsertAllDoc(collectionName, allProducts)
	if searchDeleteAllErr != nil {
		return errors.New("Fail to delete all docs")
	}

	if searchInsertAll != nil {
		log.Println(searchInsertAll)
		return errors.New("Fail to update search db with all products")
	}

	return nil

}