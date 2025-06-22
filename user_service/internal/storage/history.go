package storage

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
	"user_service/internal/models"
)

func (r *repoUser) SelectHistory(ctx context.Context, userID string, limit int) ([]models.History, error) {
	var history []models.History
	err := r.db.WithContext(ctx).Where("user_id=?", userID).Order("viewed_at DESC").Limit(limit).Find(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (r *repoUser) DeleteViewHistory(ctx context.Context, userID, productID string) error {
	var his models.History
	if productID == "" {
		err := r.db.WithContext(ctx).Where("user_id=?", userID).Delete(&his).Error
		if err != nil {
			return err
		}
		return nil
	}
	err := r.db.WithContext(ctx).Where("user_id=? and product_id=?", userID, productID).Delete(&his).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repoUser) CreateView(ctx context.Context, userID, productID uint64) error {
	view := models.History{
		UserID:    userID,
		ProductID: productID,
		ViewedAt:  time.Now(),
	}

	err := r.db.WithContext(ctx).Create(&view).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repoUser) SelectView(ctx context.Context, userID, productID uint64) (bool, error) {
	now := time.Now()
	startOfDay := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	var hist models.History
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ? AND viewed_at >= ?", userID, productID, startOfDay).
		First(&hist).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repoUser) CountView(ctx context.Context, data []uint64) (int, error) {
	var count int
	err := r.db.WithContext(ctx).Model(models.History{}).Select("COUNT(*) as count").Where("product_id in ?", data).Find(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repoUser) GetProductViewsStats(ctx context.Context, productsID []uint64) ([]models.ProductViewStat, error) {
	var rows []models.ProductViewStat
	err := r.db.WithContext(ctx).
		Model(&models.History{}).
		Select("product_id, COUNT(*) AS views").
		Where("product_id IN ?", productsID).
		Group("product_id").
		Scan(&rows).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return rows, nil
}
