package storage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"review_service/internal/models"
)

type ReviewRepository interface {
	GetReviews(ctx context.Context, data []uint64) ([]models.Review, error)
	CreateReview(ctx context.Context, review models.Review) error
	GetProductReviewsStatsHandler(ctx context.Context, productIDs []uint64) ([]models.ProductReviewStat, error)
}

type repoReview struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &repoReview{db: db}
}

func (r *repoReview) GetReviews(ctx context.Context, data []uint64) ([]models.Review, error) {
	var reviews []models.Review

	err := r.db.WithContext(ctx).Where("product_id in ?", data).Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *repoReview) CreateReview(ctx context.Context, review models.Review) error {
	err := r.db.WithContext(ctx).Save(&review).Error
	return err
}

func (r *repoReview) GetProductReviewsStatsHandler(ctx context.Context, productIDs []uint64) ([]models.ProductReviewStat, error) {
	var rows []models.ProductReviewStat
	err := r.db.WithContext(ctx).
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
