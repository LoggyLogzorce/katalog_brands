package api

import (
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func FavoriteHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, body, err := proxyTo(c, "http://localhost:8082", "", nil)
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

	productsID := models.ProductRequest{
		Favorite: favoritesID,
	}

	favoritesIDJson, err := json.Marshal(productsID)
	if err != nil {
		log.Println("FavoriteHandler: ошибка маршализации productsID:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	status, _, body, err = proxyTo(c, "http://localhost:8083", "/api/v1/products", bytes.NewBuffer(favoritesIDJson))
	if err != nil {
		log.Println("FavoriteHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("FavoriteHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить информацию о продуктах"})
		return
	}

	var products models.ProductResponse
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("FavoriteHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	c.JSON(status, products.Favorite)
}
