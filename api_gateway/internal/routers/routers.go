package routers

import (
	"api_gateway/internal/handlers"
	"api_gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Проксирование запроса к микросервису
func proxyTo(c *gin.Context, target string) {
	client := &http.Client{}
	req, _ := http.NewRequest(c.Request.Method, target+c.Request.RequestURI, c.Request.Body)
	req.Header = c.Request.Header
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "сервис недоступен"})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func SetStaticRouters(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("web/template/*")
	r.Static("/static", "./web/static")

	r.GET("/", middleware.OptionalAuthMiddleware(), handlers.HomePage)
	r.GET("/brands", middleware.OptionalAuthMiddleware(), handlers.BrandsHandler)

	return r
}

func SetApiRouters(r *gin.Engine) *gin.Engine {
	// public: просмотр каталога
	apiPublicGroup := r.Group("/api")
	{
		apiPublicGroup.GET("/brands", handlers.GetBrands)
		//apiPublicGroup.GET("/brands", func(c *gin.Context) {
		//	proxyTo(c, "http://localhost:8082")
		//})
	}
	r.GET("/brands/*any", middleware.AuthMiddleware([]string{}), func(c *gin.Context) {
		proxyTo(c, "http://catalog-service:8080")
	})
	r.GET("/products/*any", middleware.AuthMiddleware([]string{}), func(c *gin.Context) {
		proxyTo(c, "http://catalog-service:8080")
	})

	// creator: CRUD брендов и товаров
	creator := r.Group("/", middleware.AuthMiddleware([]string{}))
	{
		creator.POST("/brands", func(c *gin.Context) {
			proxyTo(c, "http://brand-service:8081")
		})
		creator.POST("/brands/:id/products", func(c *gin.Context) {
			proxyTo(c, "http://product-service:8082")
		})
	}

	// admin: модерация
	admin := r.Group("/admin", middleware.AuthMiddleware([]string{}))
	{
		admin.PUT("/brands/:id/moderate", func(c *gin.Context) {
			proxyTo(c, "http://brand-service:8081")
		})
		admin.PUT("/products/:id/moderate", func(c *gin.Context) {
			proxyTo(c, "http://product-service:8082")
		})
	}

	return r
}
