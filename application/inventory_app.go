package application

import (
	"context"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/inventory_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type InventoryApp struct {
	p *base.Persistence
	c *context.Context
}

func NewInventoryApplication(p *base.Persistence, c *context.Context) inventory_repository.InventoryHandlerRepository {
	return &InventoryApp{p, c}
}

// func (a *InventoryApp) SaveInventory(inventory *inventory_entity.Inventory) (*inventory_entity.Inventory, map[string]string) {
// 	repoInventory := inventories.NewInventoryRepository(a.p, a.c)
// 	return repoInventory.SaveInventory(inventory)
// }


func (a *InventoryApp) GetInventory(InventoryId int64) (*inventory_entity.Inventory, error) {
	repoInventory := inventories.NewInventoryRepository(a.p, a.c)
	return repoInventory.GetInventory(InventoryId)
}

func (a *InventoryApp) UpdateInventory(Inventory *inventory_entity.Inventory) (*inventory_entity.Inventory, error) {
	repoInventory := inventories.NewInventoryRepository(a.p, a.c)
	return repoInventory.UpdateInventory(Inventory)
}

func (a *InventoryApp) DeleteInventory(InventoryId int64) error {
	repoInventory := inventories.NewInventoryRepository(a.p, a.c)
	return repoInventory.DeleteInventory(InventoryId)
}

func (a *InventoryApp) GetAllInventoryInWarehouse(warehouseId int64) ([]inventory_entity.Inventory, error) {
	repoInventory := inventories.NewInventoryRepository(a.p, a.c)
	return repoInventory.GetAllInventoryInWarehouse(warehouseId)
}



