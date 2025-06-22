package storage

import (
	"product_service/internal/db"
	"product_service/internal/models"
)

func SelectCategories(limit int) ([]models.Category, error) {
	var categories []models.Category
	err := db.DB().
		Model(&models.Category{}).
		Select(`categories.*, 
        (SELECT COUNT(*) FROM products 
         WHERE products.category_id = categories.id and products.status='approved') AS product_count`).
		Limit(limit).
		Find(&categories).
		Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategory(id uint64) (models.Category, error) {
	var category models.Category
	err := db.DB().Where("id=?", id).First(&category).Error
	return category, err
}

func CreateCategory(category models.Category) (models.Category, error) {
	err := db.DB().Omit("product_count").Create(&category).Error
	return category, err
}

func UpdateCategory(category models.Category) error {
	if category.Photo == "" {
		err := db.DB().Omit("product_count", "photo").Save(&category).Error
		return err
	}
	err := db.DB().Omit("product_count").Save(&category).Error
	return err
}

func DeleteCategory(cId string) error {
	err := db.DB().Model(models.Category{}).Where("id=?", cId).Delete(nil).Error
	return err
}
