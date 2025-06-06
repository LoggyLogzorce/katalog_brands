package storage

import (
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
