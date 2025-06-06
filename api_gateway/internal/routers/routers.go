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
		apiGroup.POST("/login", api.LoginHandler)
		apiGroup.POST("/register", api.RegisterHandler)
		apiGroup.GET("/profile", middleware.OptionalAuthMiddleware(), api.ProfileHandler)
		apiGroup.GET("/favorites", middleware.OptionalAuthMiddleware(), api.FavoriteHandler)
		apiGroup.GET("/view-history", middleware.OptionalAuthMiddleware(), api.ViewHistoryHandler)
		apiGroup.PUT("/update_role", middleware.OptionalAuthMiddleware(), api.UpdateRoleHandler)
	}

	return r
}
