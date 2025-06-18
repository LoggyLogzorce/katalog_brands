package public

import (
	"api_gateway/internal/api"
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

	status, _, body, err := api.ProxyTo(c, "http://localhost:8082", "", nil)
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

	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "/api/v1/products/approved", bytes.NewBuffer(viewHistoryIDJson))
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

	var brandsID []uint64
	for _, v := range products.ViewHistory {
		brandsID = append(brandsID, v.BrandID)
	}

	brandsStruct := models.BrandRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(brandsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8084", "/api/v1/brand", bytes.NewReader(brandsIDJson))
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

	productsStruct := models.ReviewsRequest{
		ProductsID: viewHistoryID,
	}

	productsIDJson, err := json.Marshal(productsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8085", "/api/v1/get-reviews", bytes.NewReader(productsIDJson))
	if err != nil {
		log.Println("ProductsHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProductsHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var reviews []models.ReviewsResponse
	if err = json.Unmarshal(body, &reviews); err != nil {
		log.Println("ProductsHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	for i := range products.ViewHistory {
		for _, v := range brands {
			if products.ViewHistory[i].BrandID == v.ID {
				products.ViewHistory[i].Brand = v
			}
		}
		var sum float64
		var cnt int
		for _, rv := range reviews {
			if products.ViewHistory[i].ID == rv.ProductID {
				sum += rv.Rating
				cnt++
			}
		}

		products.ViewHistory[i].Rating.CountReview = cnt
		if cnt > 0 {
			products.ViewHistory[i].Rating.AvgRating = sum / float64(cnt)
		} else {
			products.ViewHistory[i].Rating.AvgRating = 0
		}
	}

	c.JSON(status, products.ViewHistory)
}

func CreateViewHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "0" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := api.ProxyTo(c, "http://localhost:8082", "", nil)
	if err != nil {
		log.Println("CreateViewHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK && status != http.StatusCreated {
		log.Println("CreateViewHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось добавить просмотр товара"})
		return
	}

	c.JSON(status, gin.H{})
}

func DeleteViewHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := api.ProxyTo(c, "http://localhost:8082", "", nil)
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

	status, _, _, err := api.ProxyTo(c, "http://localhost:8082", "", nil)
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
