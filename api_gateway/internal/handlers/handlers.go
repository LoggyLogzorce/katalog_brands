package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Brand struct {
	ID            uint64
	Name          string
	AverageRating float64
}

func HomePage(c *gin.Context) {
	role, _ := c.Get("role")
	c.HTML(200, "index.html", gin.H{
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

func GetBrands(c *gin.Context) {
	brands := []Brand{
		{1, "Dior", 4.5},
		{2, "Letuale", 4.3},
		{3, "Cosmostars", 5.0},
	}

	c.JSON(200, brands)
}
