package api

import (
	"github.com/gin-gonic/gin"
	"user_service/internal/storage"
)

func GetViewHistory(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	history, err := storage.SelectHistory(userID, -1)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения истории просмотра", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{
		"view_history": history,
	})
}
