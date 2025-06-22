package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"user_service/internal/storage"
)

type ReviewHandler struct {
	Repo storage.UserRepository
}

func NewReviewHandler(userRepo storage.UserRepository) *ReviewHandler {
	return &ReviewHandler{Repo: userRepo}
}

type UserDataRequest struct {
	UsersID []uint64 `json:"users_id"`
}

func (h *ReviewHandler) GetUserDataHandler(c *gin.Context) {
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

	usersData, err := h.Repo.SelectUsers(c.Request.Context(), data.UsersID, limit)
	if err != nil {
		log.Println("GetUserData: ошибка получения данных о пользователях", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения данных о пользователях"})
		return
	}

	c.JSON(200, usersData)
}

func (h *ReviewHandler) GetUsersHandler(c *gin.Context) {
	users, err := h.Repo.GetUsers(c.Request.Context())
	if err != nil {
		log.Println("GetUsers: ошибка получения данных о пользователях", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения данных о пользователях"})
		return
	}

	c.JSON(200, users)
}
