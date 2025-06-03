package models

import "time"

type Favorite struct {
	UserID    uint64    `json:"user_id"`
	ProductID uint64    `json:"product_id"`
	AddedAt   time.Time `json:"added_at"`
}

func (_ *Favorite) TableName() string {
	return "favorites"
}
