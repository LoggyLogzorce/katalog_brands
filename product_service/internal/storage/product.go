package storage

import (
	"gorm.io/gorm"
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

func GetProductCountsByBrand(brandIDs []uint64, role string) (map[uint64]int, error) {
	type row struct {
		BrandID uint64 `gorm:"column:brand_id"`
		Count   int    `gorm:"column:count"`
	}
	var rows []row

	if role == "creator" {
		err := db.DB().
			Model(&models.Product{}).
			Select("brand_id, COUNT(*) AS count").
			Where("brand_id IN ?", brandIDs).
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

	if status == "admin" {
		err := db.DB().
			Preload("ProductUrls").
			Preload("Category").
			Limit(limit).
			Find(&products).
			Error
		if err != nil {
			return nil, err
		}

		return products, nil
	}

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

	return products, nil
}

func GetProduct(productID, brandID, status string) (models.Product, error) {
	var product models.Product

	if status == "creator" {
		err := db.DB().
			Preload("ProductUrls").
			Preload("Category").
			Where("id=? and brand_id=?", productID, brandID).
			First(&product).
			Error
		if err != nil {
			return models.Product{}, err
		}

		return product, nil
	}

	err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Where("id=? and brand_id=? and status=?", productID, brandID, status).
		First(&product).
		Error
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func GetProductsInBrands(brandsID []uint64) ([]models.Product, error) {
	var products []models.Product

	err := db.DB().
		Preload("ProductUrls").
		Preload("Category").
		Where("brand_id in ?", brandsID).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func CreateProduct(data models.Product, urls []models.ProductUrls) (models.Product, error) {
	product := models.Product{
		BrandID:     data.BrandID,
		CategoryID:  data.CategoryID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
	}
	err := db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		for i, _ := range urls {
			urls[i].ProductID = product.ID
		}

		if err := db.DB().Create(&urls).Error; err != nil {
			return err
		}

		return nil
	})

	return product, err
}

func DeleteProduct(brandID, productID, status string) error {
	err := db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(&models.ProductUrls{}).
			Where("product_id = ?", productID).
			Delete(nil).
			Error; err != nil {
			return err
		}

		if status == "admin" {
			if err := tx.
				Model(&models.Product{}).
				Where("id = ?", productID).
				Delete(nil).
				Error; err != nil {
				return err
			}
		} else {
			if err := tx.
				Model(&models.Product{}).
				Where("id = ? AND brand_id = ?", productID, brandID).
				Delete(nil).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func UpdateProduct(data models.Product) error {
	product := models.Product{
		ID:          data.ID,
		BrandID:     data.BrandID,
		CategoryID:  data.CategoryID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
	}
	err := db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("created_at").Save(&product).Error; err != nil {
			return err
		}

		if data.ProductUrls == nil {
			return nil
		}

		if err := tx.
			Model(&models.ProductUrls{}).
			Where("product_id = ?", data.ID).
			Delete(nil).
			Error; err != nil {
			return err
		}

		if err := db.DB().Create(&data.ProductUrls).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
