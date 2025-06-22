package storage

import (
	"brand_service/internal/db"
	"brand_service/internal/models"
)

func GetAllBrandsBoot() ([]models.Brand, error) {
	var brands []models.Brand
	err := db.DB().Find(&brands).Error
	return brands, err
}
