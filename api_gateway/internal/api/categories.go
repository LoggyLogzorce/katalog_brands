package api

import (
	"api_gateway/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CategoryHandler(c *gin.Context) {
	status, _, body, err := proxyTo(c, "http://localhost:8084")
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

	var categories []models.Categories
	if err = json.Unmarshal(body, &categories); err != nil {
		log.Println("CategoryHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	c.JSON(status, gin.H{"brands": categories})
}
