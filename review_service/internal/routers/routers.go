package routers

import (
	"github.com/gin-gonic/gin"
	"review_service/internal/api"
)

func SetRouters(h *api.ReviewHandler) *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api/v1")

	apiGroup.GET("/get-reviews", h.GetReviewsHandler)
	apiGroup.POST("/create-review/:pID", h.CreateReviewHandler)
	apiGroup.GET("product_reviews_stats", h.GetProductReviewsStatsHandler)

	return r
}
