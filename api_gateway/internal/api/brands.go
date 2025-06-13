package api

import (
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func BrandsHandler(c *gin.Context) {
	count := c.Query("count")
	url := fmt.Sprintf("/api/v1/brands/approved?%s", count)
	status, _, body, err := proxyTo(c, "http://localhost:8084", url, nil)
	if err != nil {
		log.Println("BrandsHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandsHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var brands []models.Brand
	if err = json.Unmarshal(body, &brands); err != nil {
		log.Println("BrandsHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	var brandsID []uint64
	for _, v := range brands {
		brandsID = append(brandsID, v.ID)
	}

	productsCount := models.BrandCountRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(productsCount)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = proxyTo(c, "http://localhost:8083", "/api/v1/brands/count-product", bytes.NewReader(brandsIDJson))
	if err != nil {
		log.Println("BrandsHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandsHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить количество товаров брендов"})
		return
	}

	var brandsCount []models.BrandCount
	if err = json.Unmarshal(body, &brandsCount); err != nil {
		log.Println("BrandsHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	for i, v := range brands {
		for _, e := range brandsCount {
			if v.ID == e.BrandID {
				brands[i].ProductCount = e.Count
			}
		}
	}

	c.JSON(status, brands)
}

func BrandHandler(c *gin.Context) {
	status, _, body, err := proxyTo(c, "http://localhost:8084", "", nil)
	if err != nil {
		log.Println("BrandHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию бренде"})
		return
	}

	var brand models.Brand
	if err = json.Unmarshal(body, &brand); err != nil {
		log.Println("BrandHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	url := fmt.Sprintf("/api/v1/brand/%v/products/approved", brand.ID)

	status, _, body, err = proxyTo(c, "http://localhost:8083", url, nil)
	if err != nil {
		log.Println("BrandHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию бренде"})
		return
	}

	var products []models.Product
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("BrandHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	for i := range products {
		products[i].Rating.AvgRating = 3.5
		products[i].Rating.CountReview = 100
	}

	brandProducts := models.BrandResponse{
		Brand:    brand,
		Products: products,
	}

	c.JSON(http.StatusOK, brandProducts)
}
