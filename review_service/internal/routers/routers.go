package routers

import (
	"github.com/gin-gonic/gin"
	"review_service/internal/api"
)

func SetRouters() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")

	apiGroup.GET("/get-reviews", api.GetReviewsHandler)
	apiGroup.POST("/create-review/:pID", api.CreateReviewHandler)

	return r
}
