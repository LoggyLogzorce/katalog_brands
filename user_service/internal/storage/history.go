package storage

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"user_service/internal/db"
	"user_service/internal/models"
)

func SelectHistory(userID string, limit int) ([]models.History, error) {
	var history []models.History
	err := db.DB().Where("user_id=?", userID).Order("viewed_at DESC").Limit(limit).Find(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}

func DeleteViewHistory(userID, productID string) error {
	var his models.History
	if productID == "" {
		err := db.DB().Where("user_id=?", userID).Delete(&his).Error
		if err != nil {
			return err
		}
		return nil
	}
	err := db.DB().Where("user_id=? and product_id=?", userID, productID).Delete(&his).Error
	if err != nil {
		return err
	}

	return nil
}

func CreateView(userID, productID uint64) error {
	view := models.History{
		UserID:    userID,
		ProductID: productID,
		ViewedAt:  time.Now(),
	}

	err := db.DB().Create(&view).Error
	if err != nil {
		return err
	}

	return nil
}

func SelectView(userID, productID uint64) (bool, error) {
	now := time.Now()
	startOfDay := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	var hist models.History
	err := db.DB().
		Where("user_id = ? AND product_id = ? AND viewed_at >= ?", userID, productID, startOfDay).
		First(&hist).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// не смотрел сегодня
			return false, nil
		}
		// какая‑то другая ошибка при запросе
		return false, err
	}
	// запись есть
	return true, nil
}
