package application

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/warehouse_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/warehouse_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/warehouses"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type warehouseApp struct {
	p *base.Persistence
}

func NewWarehouseApplication(p *base.Persistence) warehouse_repository.WarehouseRepository {
	return &warehouseApp{p}
}


func (a *warehouseApp) SaveWarehouse(warehouse *warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, map[string]string) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p)
	return repowarehouse.SaveWarehouse(warehouse)
}

func (a *warehouseApp) GetWarehouse(warehouseId int64) (*warehouse_entity.Warehouse, error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p)
	return repowarehouse.GetWarehouse(warehouseId)
}

func (a *warehouseApp) GetAllWarehouses() ([]warehouse_entity.Warehouse, error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p)
	return repowarehouse.GetAllWarehouses()
}
	
func (a *warehouseApp) UpdateWarehouse(warehouse *warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p)
	return repowarehouse.UpdateWarehouse(warehouse)
}

func (a *warehouseApp) DeleteWarehouse(warehouseId int64) error {
	repowarehouse := warehouses.NewWareHouseRepository(a.p)
	return repowarehouse.DeleteWarehouse(warehouseId)
}

func (a *warehouseApp) SearchWarehouse(name string) ([]warehouse_entity.Warehouse, error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p)
	return repowarehouse.SearchWarehouse(name)
}

func (a *warehouseApp) UpdateWarehousesInSearchDB() (error) {
	repowarehouse := warehouses.NewWareHouseRepository(a.p)
	return repowarehouse.UpdateWarehousesInSearchDB()
}

func (a *warehouseApp) GetAllInventoryInWarehouse(warehouseId int64) (error) {
	repoInventory := inventories.NewInventoryRepository(a.p)
	_, err := repoInventory.GetAllInventoryInWarehouse(warehouseId)
	return err
}


