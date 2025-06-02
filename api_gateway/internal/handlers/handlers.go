package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomePage(c *gin.Context) {
	role, _ := c.Get("role")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Главная",
		"Role":  role,
	})
}

func BrandsHandler(c *gin.Context) {
	role, _ := c.Get("role")
	c.HTML(http.StatusOK, "brands.html", gin.H{
		"title": "Бренды",
		"Role":  role,
	})
}

func CategoriesHandler(c *gin.Context) {
	role, _ := c.Get("role")
	c.HTML(http.StatusOK, "categories.html", gin.H{
		"title": "Категории",
		"Role":  role,
	})
}

func AuthHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"title": "Вход",
	})
}

func RegisterHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Регистрация",
	})
}
