package storage

import (
	"review_service/internal/db"
	"review_service/internal/models"
)

func GetReviews(data []uint64) ([]models.Review, error) {
	var reviews []models.Review

	err := db.DB().Where("product_id in ?", data).Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func CreateReview(review models.Review) error {
	err := db.DB().Save(&review).Error
	return err
}
