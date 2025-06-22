package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"user_service/internal/models"
)

func (h *ReviewHandler) GetProfileInfoHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	role := c.GetHeader("X-Role")

	user, err := h.Repo.SelectUser(c.Request.Context(), userID, role)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных пользователя"})
		return
	}

	favorites, err := h.Repo.SelectFavorite(c.Request.Context(), userID, 6)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения избранного", "error_sys": err})
		return
	}

	history, err := h.Repo.SelectHistory(c.Request.Context(), userID, 6)
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

func (h *ReviewHandler) UpdateRoleHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		log.Println("UpdateRoleHandler: не удалось получить роль из тела запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить роль из тела запроса"})
		return
	}

	userID := c.Query("id")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}

	if err := h.Repo.UpdateRoleUser(c.Request.Context(), userID, user.Role); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Ошибка обновления роли", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})
}
