package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"product_service/internal/models"
	"product_service/internal/storage"
	"strconv"
)

func GetCategoriesHandler(c *gin.Context) {
	limit := c.Query("count")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = -1
	}

	categories, err := storage.SelectCategories(limitInt)
	if err != nil {
		log.Println("GetCategoriesHandler: ошибка получения списка категорий", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения списка категорий"})
	}

	c.JSON(200, categories)
}

func CreateCategoryHandler(c *gin.Context) {
	var data models.Category
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("CreateCategoryHandler: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	category, err := storage.CreateCategory(data)
	if err != nil {
		log.Println("CreateCategoryHandler: не удалось сохранить категорию", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось сохранить категорию"})
		return
	}

	go CreateUpdateIndexCategory(category)

	c.JSON(201, gin.H{})
}

func UpdateCategoryHandler(c *gin.Context) {
	var data models.Category
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("UpdateCategoryHandler: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	err := storage.UpdateCategory(data)
	if err != nil {
		log.Println("UpdateCategoryHandler: не удалось сохранить категорию", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось сохранить категорию"})
		return
	}

	go func() {
		category, err := storage.GetCategory(data.ID)
		if err != nil {
			log.Println("UpdateCategoryHandler: не удалось найти категорию для индексации", err)
			return
		}
		CreateUpdateIndexCategory(category)
	}()

	c.JSON(200, gin.H{})
}

func DeleteCategoryHandler(c *gin.Context) {
	categoryID := c.Param("id")

	err := storage.DeleteCategory(categoryID)
	if err != nil {
		log.Println("DeleteCategoryHandler: не удалось удалить категорию", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось удалить категорию"})
		return
	}

	go DeleteIndexCategory(categoryID)

	c.JSON(200, gin.H{})
}
