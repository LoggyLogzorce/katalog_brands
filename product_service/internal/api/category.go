package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"product_service/internal/storage"
)

func GetCategories(c *gin.Context) {
	categories, err := storage.SelectCategories()
	if err != nil {
		log.Println("GetCategories: ошибка получения списка категорий", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения списка категорий"})
	}

	c.JSON(200, categories)
}
