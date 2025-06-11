package models

import "time"

type Brand struct {
	ID          uint64    `json:"id"`
	CreatorID   uint64    `json:"creator_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LogoUrl     string    `json:"logo_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (_ *Brand) TableName() string {
	return "brands"
}
