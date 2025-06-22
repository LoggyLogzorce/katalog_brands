package routers

import (
	"auth_service/internal/api"
	"github.com/gin-gonic/gin"
)

func SetRouters(auth *api.AuthHandler) *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.POST("/register", auth.Register)
		apiGroup.POST("/login", auth.Login)
		apiGroup.GET("/validate", auth.Validate)
	}

	return r
}
