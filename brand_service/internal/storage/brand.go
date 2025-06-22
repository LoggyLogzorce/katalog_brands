package storage

import (
	"brand_service/internal/db"
	"brand_service/internal/models"
)

func GetAllBrands(status, creatorID string, limit int) ([]models.Brand, error) {
	var brands []models.Brand

	if status == "all" {
		err := db.DB().Limit(limit).Find(&brands).Error
		if err != nil {
			return nil, err
		}
		return brands, nil
	}

	if status == "creator" {
		err := db.DB().Where("creator_id = ?", creatorID).Find(&brands).Error
		if err != nil {
			return nil, err
		}
		return brands, nil
	}

	err := db.DB().Where("status=?", status).Limit(limit).Find(&brands).Error
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

func GetBrandByName(name, creatorID string) (models.Brand, error) {
	var brand models.Brand

	if creatorID != "" && creatorID != "0" {
		err := db.DB().Where("name=? and creator_id=?", name, creatorID).First(&brand).Error
		if err != nil {
			return models.Brand{}, err
		}

		return brand, nil
	}

	err := db.DB().Where("name=?", name).First(&brand).Error
	if err != nil {
		return models.Brand{}, err
	}

	return brand, nil
}

func GetBrandById(id string) (models.Brand, error) {
	var brand models.Brand

	err := db.DB().Where("id=?", id).First(&brand).Error
	if err != nil {
		return models.Brand{}, err
	}

	return brand, nil
}

func UpdateBrandInfo(brand models.Brand) (models.Brand, error) {
	err := db.DB().Save(&brand).Error
	return brand, err
}

func CreateBrand(brand models.Brand) (models.Brand, error) {
	err := db.DB().Create(&brand).Error
	return brand, err
}

func DeleteBrand(brandId string) error {
	err := db.DB().Model(models.Brand{}).Where("id=?", brandId).Delete(nil).Error
	return err
}
