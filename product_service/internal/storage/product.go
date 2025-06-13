package storage

import (
	"product_service/internal/db"
	"product_service/internal/models"
)

func SelectProduct(data []uint64, status string) ([]models.Product, error) {
	var products []models.Product

	if err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Where("status=?", status).
		Find(&products, data).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func SelectProductsInCategory(categoryID, productStatus string) ([]models.Product, error) {
	var products []models.Product

	if productStatus == "all" {
		err := db.DB().
			Preload("ProductUrls").
			Preload("Category").
			Where("category_id=?", categoryID).
			Find(&products).Error
		if err != nil {
			return nil, err
		}
		return products, nil
	}

	err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Where("category_id=? and status=?", categoryID, productStatus).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func SelectProductsInBrand(brandID, productStatus string) ([]models.Product, error) {
	var products []models.Product

	if productStatus == "all" {
		err := db.DB().
			Preload("ProductUrls").
			Preload("Category").
			Where("brand_id=?", brandID).
			Find(&products).Error
		if err != nil {
			return nil, err
		}
		return products, nil
	}

	err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Where("brand_id=? and status=?", brandID, productStatus).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductCountsByBrand(brandIDs []uint64) (map[uint64]int, error) {
	type row struct {
		BrandID uint64 `gorm:"column:brand_id"`
		Count   int    `gorm:"column:count"`
	}
	var rows []row

	err := db.DB().
		Model(&models.Product{}).
		Select("brand_id, COUNT(*) AS count").
		Where("brand_id IN ? and status='approved'", brandIDs).
		Group("brand_id").
		Scan(&rows).
		Error
	if err != nil {
		return nil, err
	}

	// переводим в map
	result := make(map[uint64]int, len(rows))
	for _, r := range rows {
		result[r.BrandID] = r.Count
	}
	return result, nil
}

func GetProducts(status string, limit int) ([]models.Product, error) {
	var products []models.Product

	err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Where("status=?", status).
		Limit(limit).
		Find(&products).
		Error
	if err != nil {
		return nil, err
	}

	return products, err
}
