package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func (h *ReviewHandler) GetFavoritesHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	favorites, err := h.Repo.SelectFavorite(c.Request.Context(), userID, -1)
	if err != nil {
		log.Println("GetFavoritesHandler: ошибка получения избранного", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения избранного", "error_sys": err})
		return
	}

	c.JSON(200, favorites)
}

func (h *ReviewHandler) CreateFavoriteHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	productID := c.Param("id")

	userIdUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Println("CreateFavoriteHandler: ошибка приведения userID к uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка приведения userID к числовому типу", "error_sys": err})
		return
	}

	productIDUint, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		log.Println("CreateFavoriteHandler: ошибка приведения productID к uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка приведения productID к числовому типу", "error_sys": err})
		return
	}

	if userIdUint == 0 || productIDUint == 0 {
		log.Printf("CreateFavoriteHandler: не подходящие данные для выполнения запроса. userID=%d productID=%d",
			userIdUint, productIDUint)
		c.AbortWithStatusJSON(400, gin.H{"error": "не подходящие данные для выполнения запроса",
			"userID": userIdUint, "productID": productID})
		return
	}

	if err = h.Repo.CreateFavorite(c.Request.Context(), userIdUint, productIDUint); err != nil {
		log.Println("CreateFavoriteHandler: ошибка добавления товара в избранное", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка добавления товара в избранное", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})

}

func (h *ReviewHandler) DeleteFavoriteHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	productID := c.Param("id")

	if err := h.Repo.DeleteFavorite(c.Request.Context(), userID, productID); err != nil {
		log.Println("DeleteFavoriteHandler: ошибка удаления товара из избранного", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка удаления товара из избранного", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})
}

func (h *ReviewHandler) ClearFavoriteHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")

	if err := h.Repo.DeleteFavorite(c.Request.Context(), userID, ""); err != nil {
		log.Println("ClearFavoriteHandler: ошибка очистки избраного", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка очистки избраного", "error_sys": err})
		return
	}

	c.JSON(200, gin.H{})
}
