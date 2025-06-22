package routers

import (
	"github.com/gin-gonic/gin"
	"product_service/internal/api"
)

func SetRouters(catHandler *api.CategoryHandler, prodHandler *api.ProductHandler) *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/products/:status", prodHandler.GetProductsHandler)
		apiGroup.GET("/categories", catHandler.GetCategoriesHandler)
		apiGroup.GET("/category/:id/products/:status", prodHandler.GetProductInCategoryHandler)
		apiGroup.POST("/category/create", catHandler.CreateCategoryHandler)
		apiGroup.PUT("/category/update/:id", catHandler.UpdateCategoryHandler)
		apiGroup.DELETE("/category/delete/:id", catHandler.DeleteCategoryHandler)
		apiGroup.GET("/brand/:id/products/:status", prodHandler.GetProductInBrandHandler)
		apiGroup.GET("/brands/products", prodHandler.GetProductInBrandsHandler)
		apiGroup.GET("/brand/:id/product/:pId", prodHandler.GetProductHandler)
		apiGroup.DELETE("/brand/:id/product/:pId", prodHandler.DeleteProductHandler)
		apiGroup.GET("/brands/count-product", prodHandler.CountProductInBrandHandler)
		apiGroup.POST("/product/create", prodHandler.CreateProductHandler)
		apiGroup.PUT("/product/:id", prodHandler.UpdateProductHandler)
	}

	return r
}
