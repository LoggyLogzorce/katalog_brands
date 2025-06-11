package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"user_service/internal/models"
	"user_service/internal/storage"
)

func GetProfileInfoHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	role := c.GetHeader("X-Role")

	user, err := storage.SelectUser(userID, role)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных пользователя"})
		return
	}

	favorites, err := storage.SelectFavorite(userID, 6)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения избранного", "error_sys": err})
		return
	}

	history, err := storage.SelectHistory(userID, 6)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения истории просмотра", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{
		"user_data":    user,
		"favorites":    favorites,
		"view_history": history,
	})
}

func UpdateRoleHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		log.Println("UpdateRoleHandler: не удалось получить роль из тела запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить роль из тела запроса"})
		return
	}

	userID := c.GetHeader("X-User-ID")
	role := c.GetHeader("X-Role")

	if err := storage.UpdateRoleUser(userID, role, user.Role); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Ошибка обновления роли", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})
}
