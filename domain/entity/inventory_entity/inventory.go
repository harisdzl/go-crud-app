package inventory_entity

type Inventory struct {
	ProductID uint64 `gorm:"primary_key;not null;" json:"productID"`
	WarehouseID uint64 `gorm:"size:100;not null;" json:"warehouseID"`
	Stock int `gorm:"size:255;not null;" json:"stock"`
}

