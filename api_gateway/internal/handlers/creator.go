package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomePageCreator(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "creator_home.html", gin.H{
		"title": "Редакция",
		"Role":  role,
		"Name":  first,
	})
}

func BrandPageCreatorHandler(c *gin.Context) {
	brandName := c.Param("name")
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "creator_brand.html", gin.H{
		"title": brandName,
		"Role":  role,
		"Name":  first,
	})
}
