package api

import (
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ViewHistoryHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, body, err := proxyTo(c, "http://localhost:8082", "", nil)
	if err != nil {
		log.Println("ViewHistoryHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ViewHistoryHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var favorites []models.ViewHistory
	if err = json.Unmarshal(body, &favorites); err != nil {
		log.Println("ViewHistoryHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	var viewHistoryID []uint64

	for _, v := range favorites {
		viewHistoryID = append(viewHistoryID, v.ProductID)
	}

	productsID := models.ProfileProductRequest{
		ViewHistory: viewHistoryID,
	}

	viewHistoryIDJson, err := json.Marshal(productsID)
	if err != nil {
		log.Println("ViewHistoryHandler: ошибка маршализации productsID:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	status, _, body, err = proxyTo(c, "http://localhost:8083", "/api/v1/products/approved", bytes.NewBuffer(viewHistoryIDJson))
	if err != nil {
		log.Println("ViewHistoryHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ViewHistoryHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить информацию о продуктах"})
		return
	}

	var products models.ProfileProductResponse
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("ViewHistoryHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	// TODO заменить на запрос к Review service
	for i := range products.ViewHistory {
		products.ViewHistory[i].Rating.AvgRating = 3.5
		products.ViewHistory[i].Rating.CountReview = 100
	}

	c.JSON(status, products.ViewHistory)
}

func CreateViewHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := proxyTo(c, "http://localhost:8082", "", nil)
	if err != nil {
		log.Println("CreateViewHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("CreateViewHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось добавить просмотр товара"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DeleteViewHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := proxyTo(c, "http://localhost:8082", "", nil)
	if err != nil {
		log.Println("DeleteViewHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("DeleteViewHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось удалить просмотр товара"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ClearViewHistoryHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := proxyTo(c, "http://localhost:8082", "", nil)
	if err != nil {
		log.Println("ClearViewHistoryHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ClearViewHistoryHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось очистить историю просмотра"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
