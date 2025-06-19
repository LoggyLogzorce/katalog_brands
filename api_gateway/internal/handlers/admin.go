package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BrandsPageAdmin(c *gin.Context) {
	role := c.GetString("role")
	c.HTML(http.StatusOK, "admin_brands.html", gin.H{
		"title": "Админ панель",
		"Role":  role,
	})
}

func ProductsPageAdmin(c *gin.Context) {
	role := c.GetString("role")
	c.HTML(http.StatusOK, "admin_products.html", gin.H{
		"title": "Админ панель",
		"Role":  role,
	})
}
