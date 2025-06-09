package models

type Rating struct {
	AvgRating   float64 `json:"avg_rating"`
	CountReview int     `json:"count_review"`
}
