package application

import (
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/inventory_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type InventoryApp struct {
	p *base.Persistence
}

func NewInventoryApplication(p *base.Persistence) inventory_repository.InventoryRepository {
	return &InventoryApp{p}
}

type inventoryAppInterface interface {
	Saveinventory(*inventory_entity.Inventory) (*inventory_entity.Inventory, map[string]string)
	GetInventory(int64) (*inventory_entity.Inventory, error)
	UpdateInventory(*inventory_entity.Inventory) (*inventory_entity.Inventory, error)
	DeleteInventory(int64) error
}

func (a *InventoryApp) SaveInventory(inventory *inventory_entity.Inventory) (*inventory_entity.Inventory, map[string]string) {
	repoInventory := inventories.NewInventoryRepository(a.p)
	return repoInventory.SaveInventory(inventory)
}


func (a *InventoryApp) GetInventory(InventoryId int64) (*inventory_entity.Inventory, error) {
	repoInventory := inventories.NewInventoryRepository(a.p)
	return repoInventory.GetInventory(InventoryId)
}

func (a *InventoryApp) UpdateInventory(Inventory *inventory_entity.Inventory) (*inventory_entity.Inventory, error) {
	repoInventory := inventories.NewInventoryRepository(a.p)
	return repoInventory.UpdateInventory(Inventory)
}

func (a *InventoryApp) DeleteInventory(InventoryId int64) error {
	repoInventory := inventories.NewInventoryRepository(a.p)
	return repoInventory.DeleteInventory(InventoryId)
}



