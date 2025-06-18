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

func FavoriteHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, body, err := api.ProxyTo(c, "http://localhost:8082", "", "", nil)
	if err != nil {
		log.Println("FavoriteHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
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

	productsID := models.ProfileProductRequest{
		Favorite: favoritesID,
	}

	favoritesIDJson, err := json.Marshal(productsID)
	if err != nil {
		log.Println("FavoriteHandler: ошибка маршализации productsID:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "", "/api/v1/products/approved", bytes.NewBuffer(favoritesIDJson))
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

	var products models.ProfileProductResponse
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("FavoriteHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	var brandsID []uint64
	for _, v := range products.Favorite {
		brandsID = append(brandsID, v.BrandID)
	}

	brandsStruct := models.BrandRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(brandsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8084", "", "/api/v1/brand", bytes.NewReader(brandsIDJson))
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
		ProductsID: favoritesID,
	}

	productsIDJson, err := json.Marshal(productsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8085", "", "/api/v1/get-reviews", bytes.NewReader(productsIDJson))
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

	for i := range products.Favorite {
		for _, v := range brands {
			if products.Favorite[i].BrandID == v.ID {
				products.Favorite[i].Brand = v
			}
		}
		var sum float64
		var cnt int
		for _, rv := range reviews {
			if products.Favorite[i].ID == rv.ProductID {
				sum += rv.Rating
				cnt++
			}
		}

		products.Favorite[i].Rating.CountReview = cnt
		if cnt > 0 {
			products.Favorite[i].Rating.AvgRating = sum / float64(cnt)
		} else {
			products.Favorite[i].Rating.AvgRating = 0
		}
	}

	c.JSON(status, products.Favorite)
}

func CreateFavoriteHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "0" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Войдите в свою учётную запись, чтобы добавлять товары в избранное"})
		return
	}
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := api.ProxyTo(c, "http://localhost:8082", "", "", nil)
	if err != nil {
		log.Println("CreateFavoriteHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("CreateFavoriteHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось добавить товар в избранное"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DeleteFavoriteHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := api.ProxyTo(c, "http://localhost:8082", "", "", nil)
	if err != nil {
		log.Println("DeleteFavoriteHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("DeleteFavoriteHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось удалить товар из избранного"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func ClearFavoriteHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := api.ProxyTo(c, "http://localhost:8082", "", "", nil)
	if err != nil {
		log.Println("ClearFavoriteHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ClearFavoriteHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось очистить избранное"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
