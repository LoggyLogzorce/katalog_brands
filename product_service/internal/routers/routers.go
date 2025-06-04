package routers

import (
	"github.com/gin-gonic/gin"
	"product_service/internal/api"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/products", api.GetProduct)
	}

	return r
}
