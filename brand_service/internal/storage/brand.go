package storage

import (
	"brand_service/internal/models"
	"context"
	"gorm.io/gorm"
)

type BrandRepository interface {
	GetAllBrands(ctx context.Context, status, creatorID string, limit int) ([]models.Brand, error)
	GetBrandInfoById(ctx context.Context, brandsID []uint64) ([]models.Brand, error)
	GetBrandByName(ctx context.Context, name, creatorID string) (models.Brand, error)
	GetBrandById(ctx context.Context, id string) (models.Brand, error)
	UpdateBrandInfo(ctx context.Context, brand models.Brand) (models.Brand, error)
	CreateBrand(ctx context.Context, brand models.Brand) (models.Brand, error)
	DeleteBrand(ctx context.Context, brandId string) error
}

type repoBrand struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &repoBrand{db: db}
}

func (r *repoBrand) GetAllBrands(ctx context.Context, status, creatorID string, limit int) ([]models.Brand, error) {
	var brands []models.Brand

	if status == "all" {
		err := r.db.WithContext(ctx).Limit(limit).Find(&brands).Error
		if err != nil {
			return nil, err
		}
		return brands, nil
	}

	if status == "creator" {
		err := r.db.WithContext(ctx).Where("creator_id = ?", creatorID).Find(&brands).Error
		if err != nil {
			return nil, err
		}
		return brands, nil
	}

	err := r.db.WithContext(ctx).Where("status=?", status).Limit(limit).Find(&brands).Error
	if err != nil {
		return nil, err
	}

	return brands, nil
}

func (r *repoBrand) GetBrandInfoById(ctx context.Context, brandsID []uint64) ([]models.Brand, error) {
	var brands []models.Brand
	err := r.db.WithContext(ctx).Find(&brands, brandsID).Error
	if err != nil {
		return nil, err
	}

	return brands, nil
}

func (r *repoBrand) GetBrandByName(ctx context.Context, name, creatorID string) (models.Brand, error) {
	var brand models.Brand

	if creatorID != "" && creatorID != "0" {
		err := r.db.WithContext(ctx).Where("name=? and creator_id=?", name, creatorID).First(&brand).Error
		if err != nil {
			return models.Brand{}, err
		}

		return brand, nil
	}

	err := r.db.WithContext(ctx).Where("name=?", name).First(&brand).Error
	if err != nil {
		return models.Brand{}, err
	}

	return brand, nil
}

func (r *repoBrand) GetBrandById(ctx context.Context, id string) (models.Brand, error) {
	var brand models.Brand

	err := r.db.WithContext(ctx).Where("id=?", id).First(&brand).Error
	if err != nil {
		return models.Brand{}, err
	}

	return brand, nil
}

func (r *repoBrand) UpdateBrandInfo(ctx context.Context, brand models.Brand) (models.Brand, error) {
	err := r.db.WithContext(ctx).Save(&brand).Error
	return brand, err
}

func (r *repoBrand) CreateBrand(ctx context.Context, brand models.Brand) (models.Brand, error) {
	err := r.db.WithContext(ctx).Create(&brand).Error
	return brand, err
}

func (r *repoBrand) DeleteBrand(ctx context.Context, brandId string) error {
	err := r.db.WithContext(ctx).Model(models.Brand{}).Where("id=?", brandId).Delete(nil).Error
	return err
}
