package models

import "time"

type ProductRequest struct {
	Favorite    []uint64 `json:"favorite"`
	ViewHistory []uint64 `json:"view_history"`
}

type ProductResponse struct {
	Favorite    []Product `json:"favorites"`
	ViewHistory []Product `json:"view_history"`
}

type Product struct {
	ID          uint64        `gorm:"primaryKey" json:"product_id"`
	BrandID     uint64        `json:"brand_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	IsFavorite  bool          `gorm:"-" json:"is_favorite"`
	ProductUrls []ProductUrls `json:"product_urls"`
	Category    Category      `gorm:"foreignKey:CategoryID;references:id" json:"category"`
}

type ProductUrls struct {
	ProductID uint64 `json:"product_id"`
	Url       string `json:"url"`
}
