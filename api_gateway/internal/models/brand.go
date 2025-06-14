package models

import "time"

type Brand struct {
	ID           uint64    `json:"id"`
	CreatorID    uint64    `json:"creator_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	LogoUrl      string    `json:"logo_url"`
	Status       string    `json:"status"`
	ProductCount int       `json:"product_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type BrandResponse struct {
	Brand    Brand     `json:"brand"`
	Products []Product `json:"products"`
}

type BrandRequest struct {
	BrandIDs []uint64 `json:"brand_ids"`
}

type BrandCount struct {
	BrandID uint64 `json:"brand_id"`
	Count   int    `json:"count"`
}
