package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"product_service/internal/es"
	"product_service/internal/models"
	"product_service/internal/storage"
	"strconv"
	"time"
)

type CategoryHandler struct {
	CatRepo storage.CategoryRepository
	EsRepo  es.IndexerRepository
}

func NewCategoryHandler(catRepo storage.CategoryRepository, esRepo es.IndexerRepository) *CategoryHandler {
	return &CategoryHandler{CatRepo: catRepo, EsRepo: esRepo}
}

func (h *CategoryHandler) GetCategoriesHandler(c *gin.Context) {
	limit := c.Query("count")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = -1
	}

	categories, err := h.CatRepo.GetAllCategories(c.Request.Context(), limitInt)
	if err != nil {
		log.Println("GetCategoriesHandler: ошибка получения списка категорий", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения списка категорий"})
	}

	c.JSON(200, categories)
}

func (h *CategoryHandler) CreateCategoryHandler(c *gin.Context) {
	var data models.Category
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("CreateCategoryHandler: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	category, err := h.CatRepo.Create(c.Request.Context(), data)
	if err != nil {
		log.Println("CreateCategoryHandler: не удалось сохранить категорию", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось сохранить категорию"})
		return
	}

	go func(c models.Category) {
		doc := es.CategoryDoc{
			ID:    fmt.Sprint(c.ID),
			Name:  c.Name,
			Photo: c.Photo,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.EsRepo.IndexCategory(ctx, doc); err != nil {
			log.Println("CreateProductHandler: ES indexing failed:", err)
			return
		}
		log.Println("CreateProductHandler: ES indexing successfully")
	}(category)

	c.JSON(201, gin.H{})
}

func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	var data models.Category
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("UpdateCategoryHandler: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	err := h.CatRepo.Update(c.Request.Context(), data)
	if err != nil {
		log.Println("UpdateCategoryHandler: не удалось сохранить категорию", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось сохранить категорию"})
		return
	}

	go func() {
		c, err := h.CatRepo.GetCategoryByID(c.Request.Context(), data.ID)
		if err != nil {
			log.Println("UpdateCategoryHandler: не удалось найти категорию для индексации", err)
			return
		}

		doc := es.CategoryDoc{
			ID:    fmt.Sprint(c.ID),
			Name:  c.Name,
			Photo: c.Photo,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.EsRepo.IndexCategory(ctx, doc); err != nil {
			log.Println("CreateProductHandler: ES indexing failed:", err)
			return
		}
		log.Println("CreateProductHandler: ES indexing successfully")
	}()

	c.JSON(200, gin.H{})
}

func (h *CategoryHandler) DeleteCategoryHandler(c *gin.Context) {
	categoryID := c.Param("id")

	err := h.CatRepo.Delete(c.Request.Context(), categoryID)
	if err != nil {
		log.Println("DeleteCategoryHandler: не удалось удалить категорию", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось удалить категорию"})
		return
	}

	go func(id string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.EsRepo.DeleteCategory(ctx, id); err != nil {
			log.Println("DeleteProductHandler: ошибка удаления из ES: doc", err)
			return
		}
		log.Println("DeleteProductHandler: успешное удаление из ES:", id)
	}(categoryID)

	c.JSON(200, gin.H{})
}
