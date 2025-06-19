package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"user_service/internal/storage"
)

type UserDataRequest struct {
	UsersID []uint64 `json:"users_id"`
}

func GetUserDataHandler(c *gin.Context) {
	var data UserDataRequest
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetUserData: ошибка разбора данных из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка разбора данных из запроса"})
		return
	}

	count := c.Query("count")
	limit, err := strconv.Atoi(count)
	if err != nil {
		limit = -1
	}

	usersData, err := storage.SelectUsers(data.UsersID, limit)
	if err != nil {
		log.Println("GetUserData: ошибка получения данных о пользователях", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения данных о пользователях"})
		return
	}

	c.JSON(200, usersData)
}

func GetUsersHandler(c *gin.Context) {
	users, err := storage.GetUsers()
	if err != nil {
		log.Println("GetUsers: ошибка получения данных о пользователях", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения данных о пользователях"})
		return
	}

	c.JSON(200, users)
}
