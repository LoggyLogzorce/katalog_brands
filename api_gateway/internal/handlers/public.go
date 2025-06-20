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

func LogoutHandler(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"localhost",
		true,
		true,
	)
	c.Redirect(http.StatusSeeOther, "/")
}

func ProfileHandler(c *gin.Context) {
	role := c.GetString("role")
	if role == "guest" {
		AuthHandler(c)
		return
	}
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

func FavoriteHandler(c *gin.Context) {
	role := c.GetString("role")
	if role == "guest" {
		AuthHandler(c)
		return
	}
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "favorite.html", gin.H{
		"title": "Избранное",
		"Role":  role,
		"Name":  first,
	})
}

func ViewHistoryHandler(c *gin.Context) {
	role := c.GetString("role")
	if role == "guest" {
		AuthHandler(c)
		return
	}
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "view_history.html", gin.H{
		"title": "История просмотра",
		"Role":  role,
		"Name":  first,
	})
}

func CategoryProductHandler(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "products_category.html", gin.H{
		"title": "Просмотр категории",
		"Role":  role,
		"Name":  first,
	})
}

func BrandPageHandler(c *gin.Context) {
	page := c.Param("name")
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "brand_page.html", gin.H{
		"title": page,
		"Role":  role,
		"Name":  first,
	})
}

func ProductsHandler(c *gin.Context) {
	page := c.Param("name")
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "products.html", gin.H{
		"title": page,
		"Role":  role,
		"Name":  first,
	})
}

func ProductHandler(c *gin.Context) {
	brand := c.Param("name")
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "product_page.html", gin.H{
		"title": "Товар" + brand,
		"Role":  role,
		"Name":  first,
	})
}

func SearchHandler(c *gin.Context) {
	role := c.GetString("role")
	name := c.GetString("name")
	var first string
	if len(name) > 0 {
		first = string([]rune(name)[0])
	}
	c.HTML(http.StatusOK, "search.html", gin.H{
		"title": "Поиск",
		"Role":  role,
		"Name":  first,
	})
}
