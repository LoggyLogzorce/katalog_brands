package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateReviewHandler(c *gin.Context) {
	userID := c.GetString("userID")
	fmt.Println(userID)
	if userID == "0" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Чтобы оставить отзыв войдите в свою учётную запись"})
		return
	}
	c.Request.Header.Set("X-User-ID", userID)

	status, _, _, err := proxyTo(c, "http://localhost:8085", "", nil)
	if err != nil {
		log.Println("CreateReview: ошибка вызова Review Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusCreated {
		log.Println("CreateReview: Review Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "Не удалось сохранить ваш отзыв, попробуйте позже"})
		return
	}

	c.JSON(status, gin.H{})
}
