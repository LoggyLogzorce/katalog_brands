package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomePage(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Главная",
		"Role":  role,
		"Name":  first,
	})
}

func BrandsHandler(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "brands.html", gin.H{
		"title": "Бренды",
		"Role":  role,
		"Name":  first,
	})
}

func CategoriesHandler(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "categories.html", gin.H{
		"title": "Категории",
		"Role":  role,
		"Name":  first,
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

func ProfileHandler(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"title": "Профиль",
		"Role":  role,
		"Name":  first,
	})
}
