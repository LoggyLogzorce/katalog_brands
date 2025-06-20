package api

import (
	"brand_service/internal/es"
	"brand_service/internal/models"
	"context"
	"fmt"
	"log"
	"time"
)

func CreateUpdateIndex(b models.Brand) {
	doc := es.BrandDoc{
		ID:          fmt.Sprint(b.ID),
		Name:        b.Name,
		Description: b.Description,
		Photo:       b.LogoUrl,
		Status:      b.Status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := indexer.IndexBrand(ctx, doc); err != nil {
		log.Println("CreateBrand: ES indexing failed:", err)
	}
	log.Println("CreateBrand: ES indexing successfully")
}

func DeleteIndex(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := indexer.DeleteBrand(ctx, id); err != nil {
		log.Println("DeleteProductHandler: ошибка удаления из ES: doc", err)
	}
	log.Println("DeleteProductHandler: успешное удаление из ES:", id)
}
