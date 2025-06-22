package routers

import (
	"github.com/gin-gonic/gin"
	"user_service/internal/api"
)

func SetRouters(h *api.ReviewHandler) *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/profile", h.GetProfileInfoHandler)
		apiGroup.GET("/user-data", h.GetUserDataHandler)
		apiGroup.GET("/get/users", h.GetUsersHandler)
		apiGroup.PUT("/update_role", h.UpdateRoleHandler)

		favGroup := apiGroup.Group("/favorites")
		{
			favGroup.GET("/", h.GetFavoritesHandler)
			favGroup.POST("/:id", h.CreateFavoriteHandler)
			favGroup.DELETE("/:id", h.DeleteFavoriteHandler)
			favGroup.DELETE("/", h.ClearFavoriteHandler)
		}

		hisGroup := apiGroup.Group("/view-history")
		{
			hisGroup.GET("/", h.GetViewHistoryHandler)
			hisGroup.POST("/:id", h.CreateViewProductHandler)
			hisGroup.DELETE("/:id", h.DeleteViewProductHandler)
			hisGroup.DELETE("/", h.ClearViewHistoryHandler)
		}
		apiGroup.GET("/product_views_stats", h.GetProductViewsStatsHandler)
	}

	return r
}
