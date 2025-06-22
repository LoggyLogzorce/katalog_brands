package storage

import (
	"product_service/internal/db"
	"product_service/internal/models"
)

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Find(&products).Error
	return products, err
}

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := db.DB().Find(&categories).Error
	return categories, err
}
