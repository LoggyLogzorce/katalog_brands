package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"user_service/internal/storage"
)

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

	if err = storage.CreateView(userIdUint, productIDUint); err != nil {
		log.Println("CreateViewProductHandler: ошибка добавления просмотра товара", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка добавления просмотра товара", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})
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
