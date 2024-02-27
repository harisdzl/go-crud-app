package categories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/category_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/category_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

// To manage new product repositories in the database

// Product Repository struct
type CategoryRepo struct {
	p *base.Persistence
	c context.Context
}

func NewCategoryRepository(p *base.Persistence, c context.Context) *CategoryRepo {
	return &CategoryRepo{p, c}
}

// To explicitly check that the ProductRepo implements the repository.ProductRepository interface
var _ category_repository.CategoryRepository = &CategoryRepo{}

func (c *CategoryRepo) SaveCategory(category *category_entity.Category) (*category_entity.Category, map[string]string) {
	cacheRepo := cache.NewCacheRepository("Redis", c.p)

	dbErr := map[string]string{}
	err := c.p.DB.Debug().Create(&category).Error
	if err != nil {
		fmt.Println("Failed to create category")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	fmt.Printf("Type of product: %T\n", category) // Log the type of the product

	cacheRepo.SetKey(fmt.Sprintf("%v_CATEGORIES", category.ID), category, time.Minute * 15)
	
	return category, nil
}

// func (c *CategoryRepo) SaveMultipleProducts(product *[]category_entity.Category) (*[]category_entity.Category, map[string]string) {
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

func (c *CategoryRepo) GetCategory(id int64) (*category_entity.Category, error) {
	var category *category_entity.Category

	cacheRepo := cache.NewCacheRepository("Redis", c.p)
	_ = cacheRepo.GetKey(fmt.Sprintf("%v_CATEGORIES", id), &category)
	if category == nil {
		err := c.p.DB.Debug().Where("id = ?", id).Take(&category).Error
		if err != nil {
			fmt.Println("Failed to get product")
			return nil, err
		}
		if category != nil && category.ID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprintf("%v_CATEGORIES", id), category, time.Minute * 15)
		}
	}


	return category, nil
}

func (c *CategoryRepo) GetAllCategories() ([]category_entity.Category, error) {
	var categories []category_entity.Category
	err := c.p.DB.Debug().Find(&categories).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return categories, nil
}

func (c *CategoryRepo) GetParentCategories(id int64) ([]category_entity.Category, error) {
	var relatedCategories []category_entity.Category
	var rootReached bool
	childCategory, err := c.GetCategory(id)
	relatedCategories = append(relatedCategories, *childCategory)
	if err != nil {
		return nil, err
	}

	// Search for direct parent category
	for !rootReached {
		parentCategory, parentCategoryErr := c.GetCategory(childCategory.ParentID)
		if parentCategoryErr != nil {
			return nil, parentCategoryErr
		}

		relatedCategories = append(relatedCategories, *parentCategory)

		if parentCategory.ParentID == 0 {
			rootReached = true
			return relatedCategories, nil
		}

		parentCategory = childCategory
	} 

	return relatedCategories, nil
}


func (c *CategoryRepo) UpdateCategory(category *category_entity.Category) (*category_entity.Category, error) {
	cacheRepo := cache.NewCacheRepository("Redis", c.p)

	err := c.p.DB.Debug().Where("id = ?", category.ID).Updates(&category).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	_ = cacheRepo.SetKey(fmt.Sprintf("%v_CATEGORIES", category.ID), category, time.Minute * 15)


	return category, nil
}

func (c *CategoryRepo) DeleteCategory(id int64) error {
	var category category_entity.Category

	err := c.p.DB.Debug().Where("id = ?", id).Delete(&category).Error
	
	cacheRepo := cache.NewCacheRepository("Redis", c.p)

	cacheRepo.DelKey(fmt.Sprintf("%v_CATEGORIES", id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}
