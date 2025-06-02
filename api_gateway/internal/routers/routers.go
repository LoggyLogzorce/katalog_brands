package routers

import (
	"api_gateway/internal/handlers"
	"api_gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"io"
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

	// Копируем все нужные заголовки из ответа сервиса (включая Authorization)
	for key, values := range resp.Header {
		for _, v := range values {
			// Gin по умолчанию не перезаписывает уже установленные заголовки, поэтому используем Add
			c.Writer.Header().Add(key, v)
		}
	}

	// Устанавливаем статус код и копируем тело
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func SetStaticRouters(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("web/template/*")
	r.Static("/static", "./web/static")

	r.GET("/", middleware.OptionalAuthMiddleware(), handlers.HomePage)
	r.GET("/auth", handlers.AuthHandler)
	r.GET("/register", handlers.RegisterHandler)
	r.GET("/brands", middleware.OptionalAuthMiddleware(), handlers.BrandsHandler)
	r.GET("/categories", middleware.OptionalAuthMiddleware(), handlers.CategoriesHandler)

	return r
}

func SetApiRouters(r *gin.Engine) *gin.Engine {
	apiPublicGroup := r.Group("/api")
	{
		apiPublicGroup.GET("/brands", func(c *gin.Context) {
			proxyTo(c, "http://localhost:8082")
		})

		apiPublicGroup.GET("/categories", func(c *gin.Context) {
			proxyTo(c, "http://localhost:8082")
		})

		apiPublicGroup.POST("/login", func(c *gin.Context) {
			proxyTo(c, "http://localhost:8081")
		})

		apiPublicGroup.POST("/register", func(c *gin.Context) {
			proxyTo(c, "http://localhost:8081")
		})
	}

	return r
}
