package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"review_service/internal/models"
	"review_service/internal/storage"
	"strconv"
	"time"
)

type ReviewHandler struct {
	Repo storage.ReviewRepository
}

func NewReviewHandler(revRepo storage.ReviewRepository) *ReviewHandler {
	return &ReviewHandler{Repo: revRepo}
}

type ProductStatsReq struct {
	ProductIDs []uint64 `json:"product_ids"`
}

type ReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewsRequest struct {
	ProductsID []uint64 `json:"products_id"`
}

func (h *ReviewHandler) GetReviewsHandler(c *gin.Context) {
	var data ReviewsRequest
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("AvgRatingHandler: ошибка разбора данных из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных из запроса"})
		return
	}

	reviews, err := h.Repo.GetReviews(c.Request.Context(), data.ProductsID)
	if err != nil {
		log.Println("AvgRatingHandler: ошибка получения отзывов о товаре", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения отзывов о товаре"})
		return
	}

	c.JSON(200, reviews)
}

func (h *ReviewHandler) CreateReviewHandler(c *gin.Context) {
	productID := c.Param("pID")
	userID := c.GetHeader("X-User-ID")
	var data ReviewRequest
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("CreateReviewHandler: ошибка разбора данных из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка разбора данных из запроса"})
		return
	}

	userIdUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Println("CreateReviewHandler: ошибка преобразования userID в uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка преобразования userID в uint64"})
		return
	}

	productIdUint, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		log.Println("CreateReviewHandler: ошибка преобразования productID в uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка преобразования productID в uint64"})
		return
	}

	review := models.Review{
		UserID:    userIdUint,
		ProductID: productIdUint,
		Rating:    data.Rating,
		Comment:   data.Comment,
		CreatedAt: time.Now(),
	}

	err = h.Repo.CreateReview(c.Request.Context(), review)
	if err != nil {
		log.Println("CreateReviewHandler: ошибка создания отзыва", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка создания отзыва"})
		return
	}

	c.JSON(201, gin.H{})
}

func (h *ReviewHandler) GetProductReviewsStatsHandler(c *gin.Context) {
	var req ProductStatsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("GetProductReviewsStatsHandler: некорректный запрос", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "некорректный запрос"})
		return
	}

	fmt.Println(req)

	rows, err := h.Repo.GetProductReviewsStatsHandler(c.Request.Context(), req.ProductIDs)
	if err != nil {
		log.Println("GetProductReviewsStatsHandler: ошибка БД", err)
		c.AbortWithStatusJSON(500, gin.H{})
		return
	}

	fmt.Println(rows)

	c.JSON(200, rows)
}
