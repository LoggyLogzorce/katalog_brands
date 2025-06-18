package routers

import (
	"api_gateway/internal/api"
	"api_gateway/internal/api/creator"
	"api_gateway/internal/api/public"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetStaticRouters(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("web/template/**/*")
	r.Static("/static", "./web/static")

	r.GET("/", middleware.OptionalAuthMiddleware(), handlers.HomePage)
	r.GET("/auth", handlers.AuthHandler)
	r.GET("/register", handlers.RegisterHandler)
	r.POST("/logout", handlers.LogoutHandler)
	r.GET("/brands", middleware.OptionalAuthMiddleware(), handlers.BrandsHandler)
	r.GET("/brand/:name", middleware.OptionalAuthMiddleware(), handlers.BrandPageHandler)
	r.GET("/brand/:name/product/:id", middleware.OptionalAuthMiddleware(), handlers.ProductHandler)
	r.GET("/categories", middleware.OptionalAuthMiddleware(), handlers.CategoriesHandler)
	r.GET("/category/:id", middleware.OptionalAuthMiddleware(), handlers.CategoryProductHandler)
	r.GET("/products", middleware.OptionalAuthMiddleware(), handlers.ProductsHandler)
	r.GET("/profile", middleware.OptionalAuthMiddleware(), handlers.ProfileHandler)
	r.GET("/favorites", middleware.OptionalAuthMiddleware(), handlers.FavoriteHandler)
	r.GET("/view-history", middleware.OptionalAuthMiddleware(), handlers.ViewHistoryHandler)

	creatorGroup := r.Group("/creator", middleware.AuthMiddleware([]string{"creator"}))
	{
		creatorGroup.GET("/brands", handlers.HomePageCreator)
		creatorGroup.GET("/brand/:name", handlers.BrandPageCreatorHandler)
	}

	r.NoRoute(middleware.OptionalAuthMiddleware(), handlers.PageNotFound)

	return r
}

func SetApiRouters(r *gin.Engine) *gin.Engine {
	publicGroup := r.Group("/api/v1")
	{
		publicGroup.GET("/brands", public.BrandsHandler)
		publicGroup.GET("/brand/:name", middleware.OptionalAuthMiddleware(), public.BrandHandler)
		publicGroup.GET("/brand/:name/product/:id", middleware.OptionalAuthMiddleware(), public.GetProductHandler)

		publicGroup.GET("/categories", public.CategoryHandler)
		publicGroup.GET("/category/:id/products/:status", middleware.OptionalAuthMiddleware(), public.CategoryProductHandler)

		publicGroup.GET("/products/:status", middleware.OptionalAuthMiddleware(), public.ProductsHandler)

		publicGroup.POST("/create-review/:pID", middleware.OptionalAuthMiddleware(), public.CreateReviewHandler)

		publicGroup.POST("/login", api.LoginHandler)
		publicGroup.POST("/register", api.RegisterHandler)
		publicGroup.GET("/profile", middleware.OptionalAuthMiddleware(), public.ProfileHandler)

		favGroup := publicGroup.Group("/favorites", middleware.OptionalAuthMiddleware())
		{
			favGroup.GET("/", public.FavoriteHandler)
			favGroup.POST("/:id", public.CreateFavoriteHandler)
			favGroup.DELETE("/:id", public.DeleteFavoriteHandler)
			favGroup.DELETE("/", public.ClearFavoriteHandler)
		}

		hisGroup := publicGroup.Group("/view-history", middleware.OptionalAuthMiddleware())
		{
			hisGroup.GET("/", public.ViewHistoryHandler)
			hisGroup.POST("/:id", public.CreateViewHandler)
			hisGroup.DELETE("/:id", public.DeleteViewHandler)
			hisGroup.DELETE("/", public.ClearViewHistoryHandler)
		}

		publicGroup.PUT("/update_role", middleware.OptionalAuthMiddleware(), public.UpdateRoleHandler)
	}

	creatorGroup := publicGroup.Group("/creator", middleware.AuthMiddleware([]string{"creator"}))
	{
		creatorGroup.GET("/brands", creator.BrandsCreatorHandler)
		creatorGroup.GET("/brand/:name", creator.BrandHandler)
		creatorGroup.PUT("/brand/:name/edit", creator.UpdateBrandHandler)
		creatorGroup.POST("/brand/:name/create-product", creator.CreateProductHandler)
		creatorGroup.POST("/brand/create", creator.CreateBrandHandler)
		creatorGroup.DELETE("/brand/:name/product/:id", creator.DeleteProductHandler)
		creatorGroup.PUT("/brand/:name/product/:id", creator.UpdateProductHandler)
	}

	return r
}
