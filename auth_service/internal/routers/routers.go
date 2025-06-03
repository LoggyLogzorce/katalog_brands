package routers

import (
	"auth_service/internal/api"
	"github.com/gin-gonic/gin"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.POST("/register", api.Register)
		apiGroup.POST("/login", api.Login)
		apiGroup.GET("/validate", api.Validate)
	}

	return r
}
