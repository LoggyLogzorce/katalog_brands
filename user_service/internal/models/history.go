package models

import "time"

type History struct {
	UserID    uint64    `json:"user_id"`
	ProductID uint64    `json:"product_id"`
	ViewedAt  time.Time `json:"viewed_at"`
}

func (_ *History) TableName() string {
	return "view_history"
}
