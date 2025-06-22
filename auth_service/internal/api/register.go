package api

import (
	"auth_service/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	if err := h.Repo.InsertUser(c.Request.Context(), req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Ошибка добавления"})
	}

	c.JSON(201, gin.H{})
}
