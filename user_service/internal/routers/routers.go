package routers

import (
	"github.com/gin-gonic/gin"
	"user_service/internal/api"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/profile", api.GetProfileInfoHandler)
		apiGroup.GET("/user-data", api.GetUserData)
		apiGroup.PUT("/update_role", api.UpdateRoleHandler)

		favGroup := apiGroup.Group("/favorites")
		{
			favGroup.GET("/", api.GetFavoritesHandler)
			favGroup.POST("/:id", api.CreateFavoriteHandler)
			favGroup.DELETE("/:id", api.DeleteFavoriteHandler)
			favGroup.DELETE("/", api.ClearFavoriteHandler)
		}

		hisGroup := apiGroup.Group("/view-history")
		{
			hisGroup.GET("/", api.GetViewHistoryHandler)
			hisGroup.POST("/:id", api.CreateViewProductHandler)
			hisGroup.DELETE("/:id", api.DeleteViewProductHandler)
			hisGroup.DELETE("/", api.ClearViewHistoryHandler)
		}
		apiGroup.GET("/product_views_stats", api.GetProductViewsStatsHandler)
	}

	return r
}
