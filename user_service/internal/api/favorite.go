package api

import (
	"github.com/gin-gonic/gin"
	"user_service/internal/storage"
)

func GetFavorites(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	favorites, err := storage.SelectFavorite(userID, -1)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения избранного", "error_sys": err})
		return
	}

	c.JSON(200, favorites)
}
