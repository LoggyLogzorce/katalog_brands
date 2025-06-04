package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	status, _, _, err := proxyTo(c, "http://localhost:8081", "", nil)
	if err != nil {
		log.Println("LoginHandler: ошибка вызова Auth Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Auth Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("LoginHandler: Auth Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось авторизоваться"})
		return
	}

	c.JSON(status, gin.H{})
}

func RegisterHandler(c *gin.Context) {
	status, _, _, err := proxyTo(c, "http://localhost:8081", "", nil)
	if err != nil {
		log.Println("RegisterHandler: ошибка вызова Auth Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Auth Service недоступен"})
		return
	}

	if status != http.StatusCreated {
		log.Println("RegisterHandler: Auth Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось зарегистрироваться"})
		return
	}

	c.JSON(status, gin.H{})
}
