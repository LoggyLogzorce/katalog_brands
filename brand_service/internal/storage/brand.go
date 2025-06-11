package storage

import (
	"brand_service/internal/db"
	"brand_service/internal/models"
)

func GetAllBrands(status string) ([]models.Brand, error) {
	var brands []models.Brand

	if status == "all" {
		err := db.DB().Find(&brands).Error
		if err != nil {
			return nil, err
		}
		return brands, nil
	}

	err := db.DB().Where("status=?", status).Find(&brands).Error
	if err != nil {
		return nil, err
	}

	return brands, nil
}

func GetBrandInfoById(brandsID []uint64) ([]models.Brand, error) {
	var brands []models.Brand
	err := db.DB().Find(&brands, brandsID).Error
	if err != nil {
		return nil, err
	}

	return brands, nil
}

func GetBrandByName(name string) (models.Brand, error) {
	var brand models.Brand
	err := db.DB().Where("name=?", name).First(&brand).Error
	if err != nil {
		return models.Brand{}, err
	}

	return brand, nil
}
