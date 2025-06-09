package storage

import (
	"time"
	"user_service/internal/db"
	"user_service/internal/models"
)

func SelectFavorite(userID string, limit int) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := db.DB().Where("user_id=?", userID).Order("added_at DESC").Limit(limit).Find(&favorites).Error
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func CreateFavorite(userID, productID uint64) error {
	favorite := models.Favorite{
		UserID:    userID,
		ProductID: productID,
		AddedAt:   time.Now(),
	}
	err := db.DB().Create(&favorite).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteFavorite(userID, productID string) error {
	var fav models.Favorite
	if productID == "" {
		err := db.DB().Where("user_id=?", userID).Delete(&fav).Error
		if err != nil {
			return err
		}
		return nil
	}
	err := db.DB().Where("user_id=? and product_id=?", userID, productID).Delete(&fav).Error
	if err != nil {
		return err
	}

	return nil
}
