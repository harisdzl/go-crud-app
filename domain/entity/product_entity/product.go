package product_entity

type Product struct {
	ID uint64 `json:"id"`
	Name string `gorm:"size:100;not null;" json:"name"`
	Description string `gorm:"size:255;not null;" json:"description"`
	Price float64 `gorm:"type:numeric;not null;" json:"price"`
	Category string `gorm:"size:100;not null;" json:"category"`
}

type ProductForInventory struct {
	ID uint64 `json:"id"`
	Name string `gorm:"size:100;not null;" json:"name"`
	Description string `gorm:"size:255;not null;" json:"description"`
	Price float64 `gorm:"type:numeric;not null;" json:"price"`
	Category string `gorm:"size:100;not null;" json:"category"`
	WarehouseID uint64 `gorm:"size:100;not null;" json:"warehouse_id"`
	Stock int `gorm:"size:255;not null;" json:"stock"`
}



