package storage

import (
	"product_service/internal/db"
	"product_service/internal/models"
)

func SelectCategories() ([]models.Category, error) {
	var categories []models.Category
	err := db.DB().
		Model(&models.Category{}).
		Select(`categories.*, 
        (SELECT COUNT(*) FROM products 
         WHERE products.category_id = categories.id) AS product_count`).
		Find(&categories).
		Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
