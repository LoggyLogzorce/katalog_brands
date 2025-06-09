package models

import "time"

type Rating struct {
	AvgRating   float64 `json:"avg_rating"`
	CountReview int     `json:"count_review"`
}

type Review struct {
	ID        uint64    `json:"id"`
	User      UserData  `json:"user"`
	Product   Product   `json:"product"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}
