package models

import "time"

type ViewHistory struct {
	UserId    uint64    `json:"user_id"`
	ProductID uint64    `json:"product_id"`
	ViewedAt  time.Time `json:"viewed_at"`
}
