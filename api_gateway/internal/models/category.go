package models

type Category struct {
	ID           uint64 `json:"category_id"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	ProductCount uint64 `json:"product_count"`
}
