package models

import "time"

type Product struct {
	ID          uint64        `gorm:"primaryKey" json:"product_id"`
	BrandID     uint64        `json:"brand_id"`
	CategoryID  uint64        `json:"category_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	IsFavorite  bool          `gorm:"-" json:"is_favorite"`
	ProductUrls []ProductUrls `json:"product_urls"`
	Category    Category      `gorm:"foreignKey:CategoryID;references:id" json:"category"`
}

func (_ *Product) TableName() string {
	return "products"
}

type ProductUrls struct {
	ProductID uint64 `json:"product_id"`
	Url       string `json:"url"`
}

func (_ *ProductUrls) TableName() string {
	return "product_photos"
}
