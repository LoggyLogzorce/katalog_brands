package admin

import (
	"api_gateway/internal/api"
	"api_gateway/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetUsersAdminHandler(c *gin.Context) {
	url := fmt.Sprintf("/api/v1/get/users")
	status, _, body, err := api.ProxyTo(c, "http://localhost:8082", "", url, nil)
	if err != nil {
		log.Println("GetUsersAdminHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("GetUsersAdminHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию для профиля"})
		return
	}

	var users []models.UserData
	if err = json.Unmarshal(body, &users); err != nil {
		log.Println("GetUsersAdminHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func CreateUserAdminHandler(c *gin.Context) {
	status, _, _, err := api.ProxyTo(c, "http://localhost:8081", "", "/api/v1/register", nil)
	if err != nil {
		log.Println("CreateUserAdminHandler: ошибка вызова Auth Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Auth Service недоступен"})
		return
	}

	if status != http.StatusCreated {
		log.Println("CreateUserAdminHandler: Auth Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию для профиля"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func UpdateRoleAdminHandler(c *gin.Context) {
	userID := c.Param("id")
	url := fmt.Sprintf("/api/v1/update_role?id=%s", userID)
	status, _, _, err := api.ProxyTo(c, "http://localhost:8082", "", url, nil)
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
