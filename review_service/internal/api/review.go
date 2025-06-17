package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"review_service/internal/models"
	"review_service/internal/storage"
	"strconv"
	"time"
)

type ReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewsRequest struct {
	ProductsID []uint64 `json:"products_id"`
}

func GetReviewsHandler(c *gin.Context) {
	var data ReviewsRequest
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("AvgRatingHandler: ошибка разбора данных из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных из запроса"})
		return
	}

	reviews, err := storage.GetReviews(data.ProductsID)
	if err != nil {
		log.Println("AvgRatingHandler: ошибка получения отзывов о товаре", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения отзывов о товаре"})
		return
	}

	c.JSON(200, reviews)
}

func CreateReviewHandler(c *gin.Context) {
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

	err = storage.CreateReview(review)
	if err != nil {
		log.Println("CreateReviewHandler: ошибка создания отзыва", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка создания отзыва"})
		return
	}

	c.JSON(201, gin.H{})
}
