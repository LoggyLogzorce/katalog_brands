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
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить информацию для профиля"})
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

	productsID := models.ProductRequest{
		Favorite:    favoritesID,
		ViewHistory: viewsHistoryID,
	}

	fmt.Println(productsID)

	productsIDJson, err := json.Marshal(productsID)
	if err != nil {
		log.Println("ProfileHandler: ошибка маршализации productsID:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	status, _, body, err = proxyTo(c, "http://localhost:8083", "/api/v1/products", bytes.NewBuffer(productsIDJson))
	if err != nil {
		log.Println("ProfileHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProfileHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить информацию о продуктах"})
		return
	}

	var products models.ProductResponse
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("ProfileHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	var resp models.ProfileResponse
	resp.UserData = profile.UserData
	resp.Favorites = products.Favorite
	resp.ViewHistory = products.ViewHistory

	c.JSON(status, resp)
}

func UpdateRoleHandler(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	c.Request.Header.Set("X-User-ID", userID)
	c.Request.Header.Set("X-Role", role)

	proxyTo(c, "http://localhost:8082", "", nil)
}
