package api

import (
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ProfileHandler(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	c.Request.Header.Set("X-User-ID", userID)
	c.Request.Header.Set("X-Role", role)

	status, _, body, err := proxyTo(c, "http://localhost:8082", "", nil)
	if err != nil {
		log.Println("ProfileHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProfileHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию для профиля"})
		return
	}

	var profile models.Profile
	if err = json.Unmarshal(body, &profile); err != nil {
		log.Println("ProfileHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	var favoritesID []uint64
	var viewsHistoryID []uint64

	for _, v := range profile.Favorites {
		favoritesID = append(favoritesID, v.ProductID)
	}

	for _, v := range profile.ViewHistory {
		viewsHistoryID = append(viewsHistoryID, v.ProductID)
	}

	productsID := models.ProfileProductRequest{
		Favorite:    favoritesID,
		ViewHistory: viewsHistoryID,
	}

	productsIDJson, err := json.Marshal(productsID)
	if err != nil {
		log.Println("ProfileHandler: ошибка маршализации productsID:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	status, _, body, err = proxyTo(c, "http://localhost:8083", "/api/v1/products/approved", bytes.NewBuffer(productsIDJson))
	if err != nil {
		log.Println("ProfileHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProfileHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию о продуктах"})
		return
	}

	var products models.ProfileProductResponse
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("ProfileHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	var brandsID []uint64
	for _, v := range products.ViewHistory {
		brandsID = append(brandsID, v.BrandID)
	}

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

	status, _, body, err = proxyTo(c, "http://localhost:8084", "/api/v1/brand", bytes.NewReader(brandsIDJson))
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

	resp := models.ProfileResponse{
		UserData:    profile.UserData,
		Favorites:   products.Favorite,
		ViewHistory: products.ViewHistory,
	}

	for i := range resp.ViewHistory {
		for _, v := range brands {
			if products.ViewHistory[i].BrandID == v.ID {
				products.ViewHistory[i].Brand = v
			}
		}
		for _, v := range favoritesID {
			if resp.ViewHistory[i].ID == v {
				resp.ViewHistory[i].IsFavorite = true
			}
		}
	}

	for i := range resp.Favorites {
		for _, v := range brands {
			if products.Favorite[i].BrandID == v.ID {
				products.Favorite[i].Brand = v
			}
		}
	}

	c.JSON(status, resp)
}

func UpdateRoleHandler(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	c.Request.Header.Set("X-User-ID", userID)
	c.Request.Header.Set("X-Role", role)

	status, _, _, err := proxyTo(c, "http://localhost:8082", "", nil)
	if err != nil {
		log.Println("UpdateRoleHandler: не удалось обновить роль пользователя", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось изменить роль пользователя"})
		return
	}

	if status != http.StatusOK {
		log.Println("UpdateRoleHandler: User service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось изменить роль пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
