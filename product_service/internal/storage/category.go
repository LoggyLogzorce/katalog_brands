package storage

import (
	"context"
	"gorm.io/gorm"
	"product_service/internal/models"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context, limit int) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, id uint64) (models.Category, error)
	Create(ctx context.Context, category models.Category) (models.Category, error)
	Update(ctx context.Context, category models.Category) error
	Delete(ctx context.Context, id string) error
}

type repoCategory struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &repoCategory{db: db}
}

func (r *repoCategory) GetAllCategories(ctx context.Context, limit int) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
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

func (r *repoCategory) GetCategoryByID(ctx context.Context, id uint64) (models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Where("id=?", id).First(&category).Error
	return category, err
}

func (r *repoCategory) Create(ctx context.Context, category models.Category) (models.Category, error) {
	err := r.db.WithContext(ctx).Omit("product_count").Create(&category).Error
	return category, err
}

func (r *repoCategory) Update(ctx context.Context, category models.Category) error {
	if category.Photo == "" {
		err := r.db.WithContext(ctx).Omit("product_count", "photo").Save(&category).Error
		return err
	}
	err := r.db.WithContext(ctx).Omit("product_count").Save(&category).Error
	return err
}

func (r *repoCategory) Delete(ctx context.Context, cId string) error {
	return r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("id = ?", cId).
		Delete(nil).
		Error
}
