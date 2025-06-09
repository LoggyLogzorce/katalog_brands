package routers

import (
	"api_gateway/internal/api"
	"api_gateway/internal/handlers"
	"api_gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetStaticRouters(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("web/template/*")
	r.Static("/static", "./web/static")

	r.GET("/", middleware.OptionalAuthMiddleware(), handlers.HomePage)
	r.GET("/auth", handlers.AuthHandler)
	r.GET("/register", handlers.RegisterHandler)
	r.GET("/brands", middleware.OptionalAuthMiddleware(), handlers.BrandsHandler)
	r.GET("/categories", middleware.OptionalAuthMiddleware(), handlers.CategoriesHandler)
	r.GET("/category/:id", middleware.OptionalAuthMiddleware(), handlers.CategoryProductHandler)
	r.GET("/profile", middleware.OptionalAuthMiddleware(), handlers.ProfileHandler)
	r.GET("/favorites", middleware.OptionalAuthMiddleware(), handlers.FavoriteHandler)
	r.GET("/view-history", middleware.OptionalAuthMiddleware(), handlers.ViewHistoryHandler)
	r.NoRoute(middleware.OptionalAuthMiddleware(), handlers.PageNotFound)

	return r
}

func SetApiRouters(r *gin.Engine) *gin.Engine {
	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/brands", api.BrandsHandler)

		apiGroup.GET("/categories", api.CategoryHandler)
		apiGroup.GET("/category/:id/products", middleware.OptionalAuthMiddleware(), api.CategoryProductHandler)

		apiGroup.POST("/login", api.LoginHandler)
		apiGroup.POST("/register", api.RegisterHandler)
		apiGroup.GET("/profile", middleware.OptionalAuthMiddleware(), api.ProfileHandler)

		favGroup := apiGroup.Group("/favorites", middleware.OptionalAuthMiddleware())
		{
			favGroup.GET("/", api.FavoriteHandler)
			favGroup.DELETE("/", api.ClearFavoriteHandler)
			favGroup.POST("/:id", api.CreateFavoriteHandler)
			favGroup.DELETE("/:id", api.DeleteFavoriteHandler)
		}

		hisGroup := apiGroup.Group("/view-history", middleware.OptionalAuthMiddleware())
		{
			hisGroup.GET("/", api.ViewHistoryHandler)
			hisGroup.POST("/:id", api.CreateViewHandler)
			hisGroup.DELETE("/:id", api.DeleteViewHandler)
			hisGroup.DELETE("/", api.ClearViewHistoryHandler)
		}

		apiGroup.PUT("/update_role", middleware.OptionalAuthMiddleware(), api.UpdateRoleHandler)
	}

	return r
}
