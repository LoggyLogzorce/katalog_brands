package models

import "time"

type Review struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	UserID    uint64    `json:"user_id"`
	ProductID uint64    `json:"product_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

func (_ *Review) TableName() string {
	return "reviews"
}

type ProductReviewStat struct {
	ProductID   uint64  `json:"product_id"`
	AvgRating   float64 `json:"avg_rating"`
	CountReview int     `json:"count_review"`
}
