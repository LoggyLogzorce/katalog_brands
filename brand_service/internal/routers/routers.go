package routers

import (
	"brand_service/internal/api"
	"github.com/gin-gonic/gin"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/brands/:status", api.GetAllBrands)
		apiGroup.GET("/brand", api.GetBrandInfo)
		apiGroup.GET("/brand/:name", api.GetBrand)
	}

	return r
}
