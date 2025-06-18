package routers

import (
	"github.com/gin-gonic/gin"
	"product_service/internal/api"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/products/:status", api.GetProductsHandler)
		apiGroup.GET("/categories", api.GetCategoriesHandler)
		apiGroup.GET("/category/:id/products/:status", api.GetProductInCategoryHandler)
		apiGroup.GET("/brand/:id/products/:status", api.GetProductInBrandHandler)
		apiGroup.GET("/brands/products", api.GetProductInBrandsHandler)
		apiGroup.GET("/brand/:id/product/:pId", api.GetProductHandler)
		apiGroup.DELETE("/brand/:id/product/:pId", api.DeleteProductHandler)
		apiGroup.GET("/brands/count-product", api.CountProductInBrandHandler)
		apiGroup.POST("/product/create", api.CreateProductHandler)
		apiGroup.PUT("/product/:id", api.UpdateProductHandler)
	}

	return r
}
