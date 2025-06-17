package models

import "time"

type Rating struct {
	AvgRating   float64 `json:"avg_rating"`
	CountReview int     `json:"count_review"`
}

type Review struct {
	ID        uint64    `json:"id"`
	User      UserData  `json:"user"`
	ProductID uint64    `json:"product"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewsRequest struct {
	ProductsID []uint64 `json:"products_id"`
}

type ReviewsResponse struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	ProductID uint64    `json:"product_id"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}
