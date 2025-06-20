package es

import (
	"brand_service/internal/storage"
	"context"
	"fmt"
	"log"
)

var (
	esClient, _ = NewClient()
	indexer     = NewIndexer(esClient)
)

func BootstrapIndexing() {
	ctx := context.Background()

	// Индексация брендов
	brands, err := storage.GetAllBrandsBoot()
	if err != nil {
		log.Printf("Ошибка при получении брендов: %v", err)
	} else {
		for _, b := range brands {
			doc := BrandDoc{
				ID:          fmt.Sprint(b.ID),
				Name:        b.Name,
				Description: b.Description,
				Photo:       b.LogoUrl,
				Status:      b.Status,
			}

			if err := indexer.IndexBrand(ctx, doc); err != nil {
				log.Printf("Ошибка при индексации бренда (ID: %d): %v", b.ID, err)
			}
		}
		log.Printf("Проиндексировано %d брендов\n", len(brands))
	}
}
