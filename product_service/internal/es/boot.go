package es

import (
	"context"
	"fmt"
	"log"
	"product_service/internal/storage"
)

var (
	esClient, _ = NewClient()
	indexer     = NewIndexer(esClient)
)

func BootstrapIndexing() {
	ctx := context.Background()

	// Индексация товаров
	products, err := storage.GetAllProducts()
	if err != nil {
		log.Printf("Ошибка при получении товаров: %v", err)
	} else {
		for _, p := range products {
			doc := ProductDoc{
				ID:          fmt.Sprint(p.ID),
				Name:        p.Name,
				Description: p.Description,
				Category:    p.Category.Name,
				Price:       p.Price,
				Photo:       p.ProductUrls[0].Url,
				Status:      p.Status,
			}
			if err := indexer.IndexProduct(ctx, doc); err != nil {
				log.Printf("Ошибка при индексации товара (ID: %d): %v", p.ID, err)
			}
		}
		log.Printf("Проиндексировано %d товаров\n", len(products))
	}

	// Индексация категорий
	categories, err := storage.GetAllCategories()
	if err != nil {
		log.Printf("Ошибка при получении категорий: %v", err)
	} else {
		for _, c := range categories {
			doc := CategoryDoc{
				ID:    fmt.Sprint(c.ID),
				Name:  c.Name,
				Photo: c.Photo,
			}
			if err := indexer.IndexCategory(ctx, doc); err != nil {
				log.Printf("Ошибка при индексации категории (ID: %d): %v", c.ID, err)
			}
		}
		log.Printf("Проиндексировано %d категорий\n", len(categories))
	}
}
