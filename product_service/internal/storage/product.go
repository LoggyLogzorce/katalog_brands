package storage

import (
	"product_service/internal/db"
	"product_service/internal/models"
)

func SelectProduct(data []uint64) ([]models.Product, error) {
	var products []models.Product

	if err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Find(&products, data).Error; err != nil {
		return nil, err
	}

	return products, nil
}
