package routers

import (
	"brand_service/internal/api"
	"github.com/gin-gonic/gin"
)

func SetRouters(h *api.BrandHandler) *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/brands/:status", h.GetAllBrands)
		apiGroup.GET("/brand", h.GetBrandInfo)
		apiGroup.GET("/brand/:name", h.GetBrand)
		apiGroup.GET("/brand/get/:id", h.GetBrandByID)
		apiGroup.PUT("/brand/:name", h.UpdateBrand)
		apiGroup.POST("/brand/create", h.CreateBrand)
		apiGroup.DELETE("/brand/:id", h.DeleteBrand)
	}

	return r
}
