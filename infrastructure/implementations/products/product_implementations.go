package products

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/product_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

// To manage new product repositories in the database

// Product Repository struct
type ProductRepo struct {
	p *base.Persistence
}

func NewProductRepository(p *base.Persistence) *ProductRepo {
	return &ProductRepo{p}
}

// To explicitly check that the ProductRepo implements the repository.ProductRepository interface
var _ product_repository.ProductRepository = &ProductRepo{}

func (r *ProductRepo) SaveProduct(product *entity.Product) (*entity.Product, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	searchRepo := search.NewSearchRepository("Mongo", r.p)

	dbErr := map[string]string{}
	err := r.p.DB.Debug().Create(&product).Error

	if err != nil {
		fmt.Println("Failed to create product")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	searchErr := searchRepo.InsertDoc(&product)

	if searchErr != nil {
		fmt.Println(searchErr)
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	cacheRepo.SetKey(fmt.Sprint(product.ID), product, time.Minute * 15)
	
	return product, nil
}

func (r *ProductRepo) SaveMultipleProducts(product *[]entity.Product) (*[]entity.Product, map[string]string) {
	dbErr := map[string]string{}
	err := r.p.DB.Debug().Create(&product).Error
	if err != nil {
		fmt.Println("Failed to create products")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	var allProducts []interface{}

    for _, p := range *product {
        allProducts = append(allProducts, p)
    }

    searchRepo := search.NewSearchRepository("Mongo", r.p)

    searchErr := searchRepo.InsertAllDoc(allProducts)

    if searchErr != nil {
        fmt.Println("Failed to insert products into search index")
        fmt.Println(searchErr)
        dbErr["search_error"] = "search index error"
        return nil, dbErr
    }



	return product, nil
}

func (r *ProductRepo) GetProduct(id int64) (*entity.Product, error) {
	var product *entity.Product

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	_ = cacheRepo.GetKey(fmt.Sprint(id), &product)
	if product == nil {
		err := r.p.DB.Debug().Where("id = ?", id).Take(&product).Error
		if err != nil {
			fmt.Println("Failed to get product")
		}
		if product != nil && product.ID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprint(id), product, time.Minute * 15)
		}
	}


	return product, nil
}

func (r *ProductRepo) GetAllProducts() ([]entity.Product, error) {
	var products []entity.Product
	err := r.p.DB.Debug().Find(&products).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepo) UpdateProduct(product *entity.Product) (*entity.Product, error) {
	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	searchRepo := search.NewSearchRepository("Mongo", r.p)

	err := r.p.DB.Debug().Where("id = ?", product.ID).Updates(&product).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	searchErr := searchRepo.UpdateDoc(uint(product.ID), &product)

	if searchErr != nil {
		return nil, err
	}

	if errors.Is(searchErr, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = cacheRepo.SetKey(fmt.Sprint(product.ID), product, time.Minute * 15)


	return product, nil
}

func (r *ProductRepo) DeleteProduct(id int64) error {
	var product entity.Product
	searchRepo := search.NewSearchRepository("Mongo", r.p)

	err := r.p.DB.Debug().Where("id = ?", id).Delete(product).Error
	searchErr := searchRepo.DeleteSingleDoc(id)
	cacheRepo := cache.NewCacheRepository("Redis", r.p)

	cacheRepo.DelKey(fmt.Sprint(id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	if searchErr != nil {
		return errors.New("database error, please try again")
	}

	return nil
}

func (r *ProductRepo) SearchProduct(name string) ([]entity.Product, error) {

	searchRepo := search.NewSearchRepository("Mongo", r.p)

    // cacheKey := fmt.Sprintf("search:%s", name)

	// Extract the results from the cursor
	var results []entity.Product

	mongoSearched, err := searchRepo.SearchDocByName(name)
	if err != nil {
		fmt.Println(err)
	}

	if err := mongoSearched.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(results) == 0 {
		fmt.Println("No such product of name: " + name)
	}

    return results, nil
}


func (r *ProductRepo) UpdateProductsInMongo() error {
	searchRepo := search.NewSearchRepository("Mongo", r.p)

	products, err := r.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var allProducts []interface{}

	for _, p := range products {
		allProducts = append(allProducts, p)
	}
	searchDeleteAllErr := searchRepo.DeleteAllDoc(allProducts)
	searchInsertAll := searchRepo.InsertAllDoc(allProducts)

	if searchDeleteAllErr != nil {
		return errors.New("Fail to delete all docs")
	}

	if searchInsertAll != nil {
		return errors.New("Fail to update mongo db with all products")
	}

	return nil

}