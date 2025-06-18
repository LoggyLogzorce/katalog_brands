package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"user_service/internal/storage"
)

type ViewCountRequest struct {
	ProductsID []uint64 `json:"products_id"`
}

type ProductStatsReq struct {
	ProductIDs []uint64 `json:"product_ids"`
}

func GetViewHistoryHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	history, err := storage.SelectHistory(userID, -1)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения истории просмотра", "error_sys": err})
		return
	}

	c.JSON(200, history)
}

func CreateViewProductHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	productID := c.Param("id")

	userIdUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Println("CreateViewProductHandler: ошибка приведения userID к uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка приведения userID к числовому типу", "error_sys": err})
		return
	}

	productIDUint, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		log.Println("CreateViewProductHandler: ошибка приведения productID к uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка приведения productID к числовому типу", "error_sys": err})
		return
	}

	if userIdUint == 0 || productIDUint == 0 {
		log.Printf("CreateViewProductHandler: не подходящие данные для выполнения запроса. userID=%d productID=%d",
			userIdUint, productIDUint)
		c.AbortWithStatusJSON(400, gin.H{"error": "не подходящие данные для выполнения запроса",
			"userID": userIdUint, "productID": productID})
		return
	}

	ok, err := storage.SelectView(userIdUint, productIDUint)
	if err != nil {
		log.Println("SelectView error:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if ok {
		// уже смотрел сегодня — ничего не делаем
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	if err = storage.CreateView(userIdUint, productIDUint); err != nil {
		log.Println("CreateViewProductHandler: ошибка добавления просмотра товара", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка добавления просмотра товара", "error_sys": err})
		return
	}

	c.JSON(201, gin.H{})
}

func DeleteViewProductHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	productID := c.Param("id")

	if err := storage.DeleteViewHistory(userID, productID); err != nil {
		log.Println("DeleteViewHistoryHandler: ошибка удаления товара из избранного", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка удаления товара из избранного", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})
}

func ClearViewHistoryHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	if err := storage.DeleteViewHistory(userID, ""); err != nil {
		log.Println("ClearFavoriteHandler: ошибка очистки избраного", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка очистки избраного", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})
}

func CountViewProductHandler(c *gin.Context) {
	var data ViewCountRequest
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("CountViewProductHandler: ошибка разбора запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка разбора запроса"})
		return
	}

	count, err := storage.CountView(data.ProductsID)
	if err != nil {
		log.Println("CountViewProductHandler: ошибка получения количества просмотров", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка получения количества просмотров"})
		return
	}

	c.JSON(200, gin.H{
		"count_view": count,
	})
}

func GetProductViewsStatsHandler(c *gin.Context) {
	var req ProductStatsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "некорректный запрос"})
		return
	}

	rows, err := storage.GetProductViewsStats(req.ProductIDs)
	if err != nil {
		log.Println("GetProductViewsStatsHandler: ошибка БД", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка БД"})
		return
	}

	c.JSON(200, rows)
}
