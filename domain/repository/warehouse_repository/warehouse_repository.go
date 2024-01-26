package warehouse_repository

import "github.com/harisquqo/quqo-challenge-1/domain/entity/warehouse_entity"

type WarehouseRepository interface {
	SaveWarehouse(*warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, map[string]string)
	GetWarehouse(int64) (*warehouse_entity.Warehouse, error)
	GetAllWarehouses() ([]warehouse_entity.Warehouse, error)
	UpdateWarehouse(*warehouse_entity.Warehouse) (*warehouse_entity.Warehouse, error)
	DeleteWarehouse(int64) error
	SearchWarehouse(string) ([]warehouse_entity.Warehouse, error)
	UpdateWarehousesInSearchDB() (error)
}