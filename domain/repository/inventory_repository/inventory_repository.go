package inventory_repository

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
)

type InventoryHandlerRepository interface {
	GetInventory(int64) (*inventory_entity.Inventory, error)
	UpdateInventory(*inventory_entity.Inventory) (*inventory_entity.Inventory, error)
}

type InventoryRepository interface {
	SaveInventory(*inventory_entity.Inventory) (*inventory_entity.Inventory, map[string]string)
	GetInventory(int64) (*inventory_entity.Inventory, error)
	UpdateInventory(*inventory_entity.Inventory) (*inventory_entity.Inventory, error)
	DeleteInventory(int64) (error)
}