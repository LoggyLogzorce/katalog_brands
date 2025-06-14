package routers

import (
	"github.com/gin-gonic/gin"
	"product_service/internal/api"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/products/:status", api.GetProducts)
		apiGroup.GET("/categories", api.GetCategories)
		apiGroup.GET("/category/:id/products/:status", api.GetProductInCategory)
		apiGroup.GET("/brand/:id/products/:status", api.GetProductInBrand)
		apiGroup.GET("/brand/:id/product/:pId", api.GetProduct)
		apiGroup.GET("/brands/count-product", api.CountProductInBrand)
	}

	return r
}
