package public

import (
	"api_gateway/internal/api"
	"api_gateway/internal/handlers"
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CategoryHandler(c *gin.Context) {
	status, _, body, err := api.ProxyTo(c, "http://localhost:8083", "GET", "", nil)
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

	status, _, body, err := api.ProxyTo(c, "http://localhost:8082", "GET", "/api/v1/favorites", nil)
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

	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "GET", "", nil)
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

	var brandsID []uint64
	for _, v := range products {
		brandsID = append(brandsID, v.BrandID)
	}

	brandsStruct := models.BrandRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(brandsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8084", "GET", "/api/v1/brand", bytes.NewReader(brandsIDJson))
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

	var productsID []uint64
	for _, v := range products {
		productsID = append(productsID, v.ID)
	}

	productsStruct := models.ReviewsRequest{
		ProductsID: productsID,
	}

	productsIDJson, err := json.Marshal(productsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8085", "GET", "/api/v1/get-reviews", bytes.NewReader(productsIDJson))
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

	for i := range products {
		for _, v := range brands {
			if products[i].BrandID == v.ID {
				products[i].Brand = v
			}
		}
		for _, v := range favorites {
			if products[i].ID == v.ProductID {
				products[i].IsFavorite = true
			}
		}
		var sum float64
		var cnt int
		for _, rv := range reviews {
			if products[i].ID == rv.ProductID {
				sum += rv.Rating
				cnt++
			}
		}

		products[i].Rating.CountReview = cnt
		if cnt > 0 {
			products[i].Rating.AvgRating = sum / float64(cnt)
		} else {
			products[i].Rating.AvgRating = 0
		}
	}

	c.JSON(status, products)
}
