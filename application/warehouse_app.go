package application

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/warehouse_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/warehouse_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/warehouses"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type warehouseApp struct {
	p *base.Persistence
	c *gin.Context
}

func NewWarehouseApplication(p *base.Persistence, c *gin.Context) warehouse_repository.WarehouseHandlerRepository {
	return &warehouseApp{p, c}
}


func (a *warehouseApp) SaveWarehouse(warehouse *warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, map[string]string) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p, a.c)
	return repowarehouse.SaveWarehouse(warehouse)
}

func (a *warehouseApp) GetWarehouse(warehouseId int64) (*warehouse_entity.Warehouse, error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p, a.c)
	return repowarehouse.GetWarehouse(warehouseId)
}

func (a *warehouseApp) GetAllWarehouses() ([]warehouse_entity.Warehouse, error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p, a.c)
	return repowarehouse.GetAllWarehouses()
}
	
func (a *warehouseApp) UpdateWarehouse(warehouse *warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p, a.c)
	return repowarehouse.UpdateWarehouse(warehouse)
}

func (a *warehouseApp) DeleteWarehouse(warehouseId int64) error {
	repowarehouse := warehouses.NewWareHouseRepository(a.p, a.c)
	return repowarehouse.DeleteWarehouse(warehouseId)
}

func (a *warehouseApp) SearchWarehouse(name string) ([]warehouse_entity.Warehouse, error) {
	searchProvider := os.Getenv("SEARCH_PROVIDER")

	searchRepo := search.NewSearchRepository(searchProvider, a.p, a.c)
	warehouseRepo := warehouses.NewWareHouseRepository(a.p, a.c)

	indexName := "warehouses"

	// Extract the results from the cursor
	var results []map[string]interface{}
	var searchProducts []warehouse_entity.Warehouse
	err := searchRepo.SearchDocByName(name, indexName, &results)

	for _, result := range results {
		warehouseId, warehouseIdErr := strconv.ParseInt(result["id"].(string), 10, 64)

		if warehouseIdErr != nil {
			return nil, warehouseIdErr
		}
		warehouse, warehouseErr := warehouseRepo.GetWarehouse(warehouseId)
		if warehouseErr != nil {
			return nil, warehouseIdErr
		}

		searchProducts = append(searchProducts, *warehouse)
	}
	if err != nil {
		fmt.Println(err)
	}

	if len(results) == 0 {
		fmt.Println("No such warehouse of name: " + name)
	}
	
	return searchProducts, nil
}

func (a *warehouseApp) UpdateWarehousesInSearchDB() (error) {
	searchProvider := os.Getenv("SEARCH_PROVIDER")
	searchRepo := search.NewSearchRepository(searchProvider, a.p, a.c)
	collectionName := "warehouses"

	warehouses, err := a.GetAllWarehouses()
	
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var allWarehouses []interface{}

    for _, p := range warehouses {
		warehouseId := fmt.Sprint(p.ID)
		searchWarehouse := map[string]interface{}{
			"id" : warehouseId,
			"name" : p.Name,
		}

        allWarehouses = append(allWarehouses, searchWarehouse)
    }

	searchDeleteAllErr := searchRepo.DeleteAllDoc(collectionName, allWarehouses)
	searchInsertAll := searchRepo.InsertAllDoc(collectionName, allWarehouses)

	if searchDeleteAllErr != nil {
		return errors.New("Fail to delete all docs")
	}

	if searchInsertAll != nil {
		return errors.New("Fail to update search db with all warehouses")
	}

	return nil
}

func (a *warehouseApp) GetAllInventoryInWarehouse(warehouseId int64) (error) {
	repoInventory := inventories.NewInventoryRepository(a.p, a.c)
	_, err := repoInventory.GetAllInventoryInWarehouse(warehouseId)
	return err
}


