package routers

import (
	"github.com/gin-gonic/gin"
	"user_service/internal/api"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/profile", api.GetProfileInfo)
		apiGroup.GET("/favorites", api.GetFavorites)
		apiGroup.GET("/view-history", api.GetViewHistory)
		apiGroup.PUT("/update_role", api.UpdateRole)
	}

	return r
}
