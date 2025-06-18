package storage

import (
	"errors"
	"gorm.io/gorm"
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

func GetProductReviewsStatsHandler(productIDs []uint64) ([]models.ProductReviewStat, error) {
	var rows []models.ProductReviewStat
	err := db.DB().
		Model(&models.Review{}).
		Select("product_id, AVG(rating) AS avg_rating, COUNT(*) AS count_review").
		Where("product_id IN ?", productIDs).
		Group("product_id").
		Scan(&rows).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return rows, nil
}
