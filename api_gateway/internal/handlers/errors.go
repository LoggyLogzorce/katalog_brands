package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PageNotFound(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"title": "Не найдено",
		"path":  c.Request.URL.Path,
		"Role":  role,
		"Name":  first,
	})
}
