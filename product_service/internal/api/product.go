package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"product_service/internal/models"
	"product_service/internal/storage"
)

func GetProduct(c *gin.Context) {
	var data models.ProductResponse
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetProduct: ошибка получения данных из запроса")
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных из запроса"})
		return
	}

	productsFavorite, err := storage.SelectProduct(data.Favorite)
	if err != nil {
		log.Println("GetProduct: ошибка получения данных о продуктах")
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
		return
	}

	productsViewHistory, err := storage.SelectProduct(data.ViewHistory)
	if err != nil {
		log.Println("GetProduct: ошибка получения данных о продуктах")
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
		return
	}

	fmt.Println(productsFavorite)

	c.JSON(200, gin.H{
		"favorites":    productsFavorite,
		"view_history": productsViewHistory,
	})

}
