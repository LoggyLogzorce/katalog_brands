package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BrandsPageAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_brands.html", gin.H{})
}

func ProductsPageAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_products.html", gin.H{})
}

func CategoriesPageAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_categories.html", gin.H{})
}

func UsersPageAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_users.html", gin.H{})
}
