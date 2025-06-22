package storage

import (
	"context"
	"gorm.io/gorm"
	"product_service/internal/models"
)

type ProductRepository interface {
	SelectProduct(ctx context.Context, data []uint64, status string) ([]models.Product, error)
	SelectProductsInCategory(ctx context.Context, categoryID, productStatus string) ([]models.Product, error)
	SelectProductsInBrand(ctx context.Context, brandID, productStatus string) ([]models.Product, error)
	GetProductCountsByBrand(ctx context.Context, brandIDs []uint64, role string) (map[uint64]int, error)
	GetProducts(ctx context.Context, status string, limit int) ([]models.Product, error)
	GetProduct(ctx context.Context, productID, brandID, status string) (models.Product, error)
	GetProductsInBrands(ctx context.Context, brandsID []uint64) ([]models.Product, error)
	CreateProduct(ctx context.Context, data models.Product, urls []models.ProductUrls) (models.Product, error)
	UpdateProduct(ctx context.Context, data models.Product) error
	Delete(ctx context.Context, brandID, productID, status string) error
}

type repoProduct struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &repoProduct{db: db}
}

func (r *repoProduct) SelectProduct(ctx context.Context, data []uint64, status string) ([]models.Product, error) {
	var products []models.Product

	if err := r.db.WithContext(ctx).
		Preload("ProductUrls").
		Preload("Category").
		Where("status=?", status).
		Find(&products, data).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repoProduct) SelectProductsInCategory(ctx context.Context, categoryID, productStatus string) ([]models.Product, error) {
	var products []models.Product

	if productStatus == "all" {
		err := r.db.WithContext(ctx).
			Preload("ProductUrls").
			Preload("Category").
			Where("category_id=?", categoryID).
			Find(&products).Error
		if err != nil {
			return nil, err
		}
		return products, nil
	}

	err := r.db.WithContext(ctx).
		Preload("ProductUrls").
		Preload("Category").
		Where("category_id=? and status=?", categoryID, productStatus).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repoProduct) SelectProductsInBrand(ctx context.Context, brandID, productStatus string) ([]models.Product, error) {
	var products []models.Product

	if productStatus == "all" {
		err := r.db.WithContext(ctx).
			Preload("ProductUrls").
			Preload("Category").
			Where("brand_id=?", brandID).
			Find(&products).Error
		if err != nil {
			return nil, err
		}
		return products, nil
	}

	err := r.db.WithContext(ctx).
		Preload("ProductUrls").
		Preload("Category").
		Where("brand_id=? and status=?", brandID, productStatus).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repoProduct) GetProductCountsByBrand(ctx context.Context, brandIDs []uint64, role string) (map[uint64]int, error) {
	type row struct {
		BrandID uint64 `gorm:"column:brand_id"`
		Count   int    `gorm:"column:count"`
	}
	var rows []row

	if role == "creator" {
		err := r.db.WithContext(ctx).
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

	err := r.db.WithContext(ctx).
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

func (r *repoProduct) GetProducts(ctx context.Context, status string, limit int) ([]models.Product, error) {
	var products []models.Product

	if status == "admin" {
		err := r.db.WithContext(ctx).
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

	err := r.db.WithContext(ctx).
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

func (r *repoProduct) GetProduct(ctx context.Context, productID, brandID, status string) (models.Product, error) {
	var product models.Product

	if status == "creator" {
		err := r.db.WithContext(ctx).
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

	err := r.db.WithContext(ctx).
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

func (r *repoProduct) GetProductsInBrands(ctx context.Context, brandsID []uint64) ([]models.Product, error) {
	var products []models.Product

	err := r.db.WithContext(ctx).
		Preload("ProductUrls").
		Preload("Category").
		Where("brand_id in ?", brandsID).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repoProduct) CreateProduct(ctx context.Context, data models.Product, urls []models.ProductUrls) (models.Product, error) {
	product := models.Product{
		BrandID:     data.BrandID,
		CategoryID:  data.CategoryID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		for i, _ := range urls {
			urls[i].ProductID = product.ID
		}

		if err := tx.Create(&urls).Error; err != nil {
			return err
		}

		return nil
	})

	return product, err
}

func (r *repoProduct) Delete(ctx context.Context, brandID, productID, status string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (r *repoProduct) UpdateProduct(ctx context.Context, data models.Product) error {
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
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		if err := tx.Create(&data.ProductUrls).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
