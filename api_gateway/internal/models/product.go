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
	ID          uint64    `json:"product_id"`
	BrandID     uint64    `json:"brand_id"`
	CategoryID  uint64    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
