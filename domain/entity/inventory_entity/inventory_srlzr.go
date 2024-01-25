package inventory_entity

import "github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"


func ConvertProductInventoryToInventory(productInventory product_entity.ProductForInventory, product product_entity.Product) Inventory {
	var inventory Inventory

	inventory.ProductID = product.ID
	inventory.WarehouseID = productInventory.WarehouseID
	inventory.Stock = productInventory.Stock

	return inventory
}