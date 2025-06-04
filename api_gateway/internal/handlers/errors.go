package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PageNotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"path": c.Request.URL.Path,
	})
}

func PageForbidden(c *gin.Context) {
	c.HTML(http.StatusForbidden, "403.html", gin.H{
		"path": c.Request.URL.Path,
	})
}
