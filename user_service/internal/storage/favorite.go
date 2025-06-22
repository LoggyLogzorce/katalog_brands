package storage

import (
	"context"
	"time"
	"user_service/internal/models"
)

func (r *repoUser) SelectFavorite(ctx context.Context, userID string, limit int) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := r.db.WithContext(ctx).Where("user_id=?", userID).Order("added_at DESC").Limit(limit).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (r *repoUser) CreateFavorite(ctx context.Context, userID, productID uint64) error {
	favorite := models.Favorite{
		UserID:    userID,
		ProductID: productID,
		AddedAt:   time.Now(),
	}
	err := r.db.WithContext(ctx).Create(&favorite).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repoUser) DeleteFavorite(ctx context.Context, userID, productID string) error {
	var fav models.Favorite
	if productID == "" {
		err := r.db.WithContext(ctx).Where("user_id=?", userID).Delete(&fav).Error
		if err != nil {
			return err
		}
		return nil
	}
	err := r.db.WithContext(ctx).Where("user_id=? and product_id=?", userID, productID).Delete(&fav).Error
	if err != nil {
		return err
	}

	return nil
}
