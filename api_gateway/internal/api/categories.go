package api

import (
	"api_gateway/internal/handlers"
	"api_gateway/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CategoryHandler(c *gin.Context) {
	status, _, body, err := proxyTo(c, "http://localhost:8083", "", nil)
	if err != nil {
		log.Println("CategoryHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("CategoryHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var categories []models.Category
	if err = json.Unmarshal(body, &categories); err != nil {
		log.Println("CategoryHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	c.JSON(status, categories)
}

func CategoryProductHandler(c *gin.Context) {
	productStatus := c.Param("status")
	if productStatus != "approved" {
		handlers.PageNotFound(c)
		return
	}

	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, body, err := proxyTo(c, "http://localhost:8082", "/api/v1/favorites", nil)
	if err != nil {
		log.Println("FavoriteHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("FavoriteHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var favorites []models.Favorite
	if err = json.Unmarshal(body, &favorites); err != nil {
		log.Println("FavoriteHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	var favoritesID []uint64

	for _, v := range favorites {
		favoritesID = append(favoritesID, v.ProductID)
	}

	status, _, body, err = proxyTo(c, "http://localhost:8083", "", nil)
	if err != nil {
		log.Println("CategoryProductHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("CategoryProductHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var products []models.Product
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("CategoryProductHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	for i := range products {
		for _, v := range favoritesID {
			if products[i].ID == v {
				products[i].IsFavorite = true
			}
		}
	}

	for i := range products {
		products[i].Rating.AvgRating = 3.5
		products[i].Rating.CountReview = 100
	}

	c.JSON(status, products)
}
