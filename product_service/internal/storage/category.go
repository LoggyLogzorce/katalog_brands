package storage

import (
	"product_service/internal/db"
	"product_service/internal/models"
)

func SelectCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := db.DB().Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
