package api

import (
	"context"
	"fmt"
	"log"
	"product_service/internal/es"
	"product_service/internal/models"
	"time"
)

var (
	esClient, _ = es.NewClient()
	indexer     = es.NewIndexer(esClient)
)

func CreateUpdateIndexProduct(p models.Product) {
	doc := es.ProductDoc{
		ID:          fmt.Sprint(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category.Name,
		Price:       p.Price,
		Photo:       p.ProductUrls[0].Url,
		Status:      p.Status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := indexer.IndexProduct(ctx, doc); err != nil {
		log.Println("CreateProductHandler: ES indexing failed:", err)
		return
	}
	log.Println("CreateProductHandler: ES indexing successfully")
}

func DeleteIndexProduct(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := indexer.DeleteProduct(ctx, id); err != nil {
		log.Println("DeleteProductHandler: ошибка удаления из ES: doc", err)
		return
	}
	log.Println("DeleteProductHandler: успешное удаление из ES:", id)
}

func CreateUpdateIndexCategory(c models.Category) {
	doc := es.CategoryDoc{
		ID:    fmt.Sprint(c.ID),
		Name:  c.Name,
		Photo: c.Photo,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := indexer.IndexCategory(ctx, doc); err != nil {
		log.Println("CreateProductHandler: ES indexing failed:", err)
		return
	}
	log.Println("CreateProductHandler: ES indexing successfully")
}

func DeleteIndexCategory(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := indexer.DeleteCategory(ctx, id); err != nil {
		log.Println("DeleteProductHandler: ошибка удаления из ES: doc", err)
		return
	}
	log.Println("DeleteProductHandler: успешное удаление из ES:", id)
}
