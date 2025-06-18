package api

import (
	"github.com/gin-gonic/gin"
	"log"
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
