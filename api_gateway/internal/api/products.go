package api

import (
	"api_gateway/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ProductsHandler(c *gin.Context) {
	count := c.Query("count")
	url := fmt.Sprintf("/api/v1/products/approved?count=%s", count)

	status, _, body, err := proxyTo(c, "http://localhost:8083", url, nil)
	if err != nil {
		log.Println("ProductsHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProductsHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var products []models.Product
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("ProductsHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	c.JSON(http.StatusOK, products)
}
