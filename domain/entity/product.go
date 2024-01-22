package entity

type Product struct {
	ID uint64 `json:"id"`
	Name string `gorm:"size:100;not null;" json:"name"`
	Description string `gorm:"size:255;not null;" json:"description"`
	Price float64 `gorm:"type:numeric;not null;" json:"price"`
	Category string `gorm:"size:100;not null;" json:"category"`
	Stock int `gorm:"type:numeric;not null;" json:"stock"`
	Image string `gorm:"not null;" json:"image"`
}

// type ProductRequestToCreate struct {
// 	Name string `json:"name" bindi`
// 	Description string `gorm:"size:255;not null;" json:"description"`
// 	Price float64 `gorm:"type:numeric;not null;" json:"price"`
// 	Category string `gorm:"size:100;not null;" json:"category"`
// 	Stock uint64 `gorm:"type:numeric;not null;" json:"stock"`
// 	Image string `gorm:"not null;" json:"image"`
// }

