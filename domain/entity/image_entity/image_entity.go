package image_entity

type ImageRaw struct {
	ID uint64 `json:"id"`
	ProductID string `gorm:"size:100;not null;" json:"productID"`
	Caption string `gorm:"not null;" json:"caption"`
}

type Image struct {
	ID uint64 `json:"id"`
	ProductID string `gorm:"size:100;not null;" json:"productID"`
	Caption string `gorm:"not null;" json:"caption"`
	Url string `gorm:"not null;" json:"Url"`
}
