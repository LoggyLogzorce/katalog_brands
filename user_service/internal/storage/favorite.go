package storage

import (
	"user_service/internal/db"
	"user_service/internal/models"
)

func SelectFavorite(userID string, limit int) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := db.DB().Where("user_id=?", userID).Select("user_id", "product_id", "added_at").Limit(limit).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	return favorites, nil
}
