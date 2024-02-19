package application

import (
	"errors"
	"fmt"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/inventory_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/inventories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type InventoryApp struct {
	p *base.Persistence
}

func NewInventoryApplication(p *base.Persistence) inventory_repository.InventoryHandlerRepository {
	return &InventoryApp{p}
}

// func (a *InventoryApp) SaveInventory(inventory *inventory_entity.Inventory) (*inventory_entity.Inventory, map[string]string) {
// 	repoInventory := inventories.NewInventoryRepository(a.p)
// 	return repoInventory.SaveInventory(inventory)
// }


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

func (a *InventoryApp) GetAllInventoryInWarehouse(warehouseId int64) ([]inventory_entity.Inventory, error) {
	repoInventory := inventories.NewInventoryRepository(a.p)
	return repoInventory.GetAllInventoryInWarehouse(warehouseId)
}

func (a *InventoryApp) ReduceInventory(productId int64, quantityOrdered int64) error {
	repoInventory := inventories.NewInventoryRepository(a.p)
	inventory, err := repoInventory.GetInventory(productId)

	if err != nil {
		return err
	}

	if inventory.Stock < int(quantityOrdered) {
		return errors.New(fmt.Sprintf("Maximum quantity is %v", inventory.Stock))
	}

	newStock := inventory.Stock - int(quantityOrdered)

	inventory.Stock = newStock
	_, updateErr := a.UpdateInventory(inventory)

	if updateErr != nil {
		return errors.New("Error in updating stock")
	}

	return nil
}


func (a *InventoryApp) IncreaseInventory(productId int64, quantityToAdd int64) error {
	repoInventory := inventories.NewInventoryRepository(a.p)
	inventory, err := repoInventory.GetInventory(productId)

	if err != nil {
		return err
	}

	newStock := inventory.Stock + int(quantityToAdd)

	inventory.Stock = newStock
	_, updateErr := a.UpdateInventory(inventory)

	if updateErr != nil {
		return errors.New("Error in updating stock")
	}
	
	return nil
}




