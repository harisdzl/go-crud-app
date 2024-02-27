package warehouses

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/warehouse_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/warehouse_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

// To manage new warehouse repositories in the database

// warehouse Repository struct
type WarehouseRepo struct {
	p *base.Persistence
	c context.Context
}

func NewWareHouseRepository(p *base.Persistence, c context.Context) *WarehouseRepo {
	return &WarehouseRepo{p, c}
}

// To explicitly check that the WarehouseRepo implements the repository.WarehouseRepository interface
var _ warehouse_repository.WarehouseRepository = &WarehouseRepo{}

func (r *WarehouseRepo) SaveWarehouse(warehouse *warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	searchRepo := search.NewSearchRepository("Mongo", r.p)

	dbErr := map[string]string{}
	err := r.p.DB.Debug().Create(&warehouse).Error
	collectionName := "warehouses"
	if err != nil {
		fmt.Println("Failed to create warehouse")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	searchErr := searchRepo.InsertDoc(collectionName, &warehouse)

	if searchErr != nil {
		fmt.Println(searchErr)
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	cacheRepo.SetKey(fmt.Sprintf("%v_WAREHOUSE", warehouse.ID), warehouse, time.Minute * 15)
	
	return warehouse, nil
}


func (r *WarehouseRepo) GetWarehouse(id int64) (*warehouse_entity.Warehouse, error) {
	var warehouse *warehouse_entity.Warehouse

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	_ = cacheRepo.GetKey(fmt.Sprintf("%v_WAREHOUSE", id), &warehouse)
	if warehouse == nil {
		err := r.p.DB.Debug().Where("id = ?", id).Take(&warehouse).Error
		if err != nil {
			fmt.Println("Failed to get warehouse")
		}
		if warehouse != nil && warehouse.ID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprintf("%v_WAREHOUSE", id), warehouse, time.Minute * 15)
		}
	}


	return warehouse, nil
}

func (r *WarehouseRepo) GetAllWarehouses() ([]warehouse_entity.Warehouse, error) {
	var warehouses []warehouse_entity.Warehouse
	err := r.p.DB.Debug().Find(&warehouses).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return warehouses, nil
}

func (r *WarehouseRepo) UpdateWarehouse(warehouse *warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, error) {
	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	searchRepo := search.NewSearchRepository("Mongo", r.p)
	collectionName := "warehouses"

	err := r.p.DB.Debug().Where("id = ?", warehouse.ID).Updates(&warehouse).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	searchErr := searchRepo.UpdateDoc(uint(warehouse.ID), collectionName, &warehouse)

	if searchErr != nil {
		return nil, err
	}

	if errors.Is(searchErr, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = cacheRepo.SetKey(fmt.Sprintf("%v_WAREHOUSE", warehouse.ID), warehouse, time.Minute * 15)


	return warehouse, nil
}

func (r *WarehouseRepo) DeleteWarehouse(id int64) error {
	var warehouse warehouse_entity.Warehouse
	searchRepo := search.NewSearchRepository("Mongo", r.p)
	collectionName := "warehouses"
	fieldName := "id"
	err := r.p.DB.Debug().Where("id = ?", id).Delete(&warehouse).Error
	
	searchErr := searchRepo.DeleteSingleDoc(fieldName, collectionName, id)
	cacheRepo := cache.NewCacheRepository("Redis", r.p)

	cacheRepo.DelKey(fmt.Sprintf("%v_WAREHOUSE", id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	if searchErr != nil {
		return errors.New("database error, please try again")
	}

	return nil
}

func (r *WarehouseRepo) SearchWarehouse(name string) ([]warehouse_entity.Warehouse, error) {

	searchRepo := search.NewSearchRepository("Mongo", r.p)
	indexName := "warehouses"

	// Extract the results from the cursor
	var results []warehouse_entity.Warehouse

	err := searchRepo.SearchDocByName(name, indexName, &results)
	
	if err != nil {
		fmt.Println(err)
	}

	if len(results) == 0 {
		fmt.Println("No such warehouse of name: " + name)
	}

    return results, nil
}


func (r *WarehouseRepo) UpdateWarehousesInSearchDB() error {
	searchRepo := search.NewSearchRepository("Mongo", r.p)
	collectionName := "warehouses"

	warehouses, err := r.GetAllWarehouses()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var allwarehouses []interface{}

	for _, p := range warehouses {
		allwarehouses = append(allwarehouses, p)
	}
	searchDeleteAllErr := searchRepo.DeleteAllDoc(collectionName, allwarehouses)
	searchInsertAll := searchRepo.InsertAllDoc(collectionName, allwarehouses)

	if searchDeleteAllErr != nil {
		return errors.New("Fail to delete all docs")
	}

	if searchInsertAll != nil {
		return errors.New("Fail to update mongo db with all warehouses")
	}

	return nil

}