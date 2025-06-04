package storage

import (
	"product_service/internal/db"
	"product_service/internal/models"
)

func SelectProduct(data []uint64) ([]models.Product, error) {
	var products []models.Product

	//for _, v := range data {
	//	arrId = append(arrId, v.ID)
	//}

	if err := db.DB().Find(&products, data).Error; err != nil {
		return nil, err
	}

	return products, nil
}
