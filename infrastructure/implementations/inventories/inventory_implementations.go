package inventories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/inventory_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"gorm.io/gorm"
)

// To manage new Inventory repositories in the database

// Inventory Repository struct
type InventoryRepo struct {
	p *base.Persistence
}

func NewInventoryRepository(p *base.Persistence) *InventoryRepo {
	return &InventoryRepo{p}
}

// To explicitly check that the InventoryRepo implements the repository.InventoryRepository interface
var _ inventory_repository.InventoryRepository = &InventoryRepo{}

func (r *InventoryRepo) SaveInventory(inventory *inventory_entity.Inventory) (*inventory_entity.Inventory, map[string]string) {

	cacheRepo := cache.NewCacheRepository("Redis", r.p)

	dbErr := map[string]string{}
	err := r.p.DB.Debug().Create(&inventory).Error

	if err != nil {
		fmt.Println("Failed to create inventory")
		fmt.Println(err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}

	cacheRepo.SetKey(fmt.Sprintf("%v_INVENTORY", inventory.ProductID), inventory, time.Minute * 15)
	
	return inventory, nil
}


func (r *InventoryRepo) GetInventory(productID int64) (*inventory_entity.Inventory, error) {
	var inventory *inventory_entity.Inventory

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	_ = cacheRepo.GetKey(fmt.Sprintf("%v_INVENTORY", productID), &inventory)
	if inventory == nil {
		err := r.p.DB.Debug().Where("product_id = ?", productID).Take(&inventory).Error
		if err != nil {
			fmt.Println("Failed to get Inventory")
		}
		if inventory != nil && inventory.ProductID > 0 {
			_ = cacheRepo.SetKey(fmt.Sprintf("%v_INVENTORY", productID), inventory, time.Minute * 15)
		}
	}


	return inventory, nil
}


func (r *InventoryRepo) GetAllInventoryInWarehouse(warehouseID int64) ([]inventory_entity.Inventory, error) {
	var inventory []inventory_entity.Inventory

	err := r.p.DB.Debug().Where("warehouse_id = ?", warehouseID).Find(&inventory).Error

	if err != nil {
		fmt.Println("Failed to get Inventory")
	}

	return inventory, nil
}


func (r *InventoryRepo) UpdateInventory(inventory *inventory_entity.Inventory) (*inventory_entity.Inventory, error) {
	cacheRepo := cache.NewCacheRepository("Redis", r.p)

	err := r.p.DB.Debug().Where("product_id = ?", inventory.ProductID).Updates(&inventory).Error
	if err != nil {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = cacheRepo.SetKey(fmt.Sprintf("%v_INVENTORY", inventory.ProductID), inventory, time.Minute * 15)

	return inventory, nil
}

func (r *InventoryRepo) DeleteInventory(id int64) error {
	var inventory inventory_entity.Inventory	

	err := r.p.DB.Debug().Where("product_id = ?", id).Delete(&inventory).Error
	if err != nil {
		return err
	}

	cacheRepo := cache.NewCacheRepository("Redis", r.p)
	


	cacheRepo.DelKey(fmt.Sprintf("%v_INVENTORY", id))
	if err != nil {
		return errors.New("database error, please try again")
	}

	return nil
}

func (r *InventoryRepo) ReduceInventory(tx *gorm.DB, id int64, quantityOrdered int64) error {
	// Update inventory stock directly in the database
	if tx == nil {
		var errTx error
		tx := r.p.DB.Begin()
		if tx.Error != nil {
			return errors.New("Failed to start transaction")
		}
	
		// Defer rollback in case of panic
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			} else if errTx != nil {
				tx.Rollback()
			} else {
				errC := tx.Commit().Error
				if errC != nil {
					tx.Rollback()
				}
			}
		}()
	}
	inventory, invErr := r.GetInventory(id)
	if invErr != nil {
		return invErr
	}

	if inventory.Stock < int(quantityOrdered) {
		return errors.New(fmt.Sprintf("Not enough stock. Maximum quantity is %v", inventory.Stock))
	}
	result := tx.Model(&inventory).
		Update("stock", gorm.Expr("stock - ?", quantityOrdered))
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errors.New("inventory not found")
	}

	// Check if the stock is negative after reduction
	if result.RowsAffected < 0 {
		log.Println("Not enough stock")
		return errors.New("not enough stock")
	}

	return nil
}

func (r *InventoryRepo) IncreaseInventory(productId int64, quantityToAdd int64) error {
    // Update inventory stock directly in the database
    result := r.p.DB.Model(&inventory_entity.Inventory{}).Where("product_id = ?", productId).
        Update("stock", gorm.Expr("stock + ?", quantityToAdd))
    if result.Error != nil {
        return result.Error
    }

    // Check if any rows were affected
    if result.RowsAffected == 0 {
        return errors.New("inventory not found")
    }

    return nil
}
