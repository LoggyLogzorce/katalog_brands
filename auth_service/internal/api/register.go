package api

import (
	"auth_service/internal/models"
	"auth_service/internal/storage"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	fmt.Println(req)

	if err := storage.InsertUser(req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Ошибка добавления"})
	}

	c.JSON(201, gin.H{})
}
