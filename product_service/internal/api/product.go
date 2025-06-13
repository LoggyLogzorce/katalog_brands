package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"product_service/internal/models"
	"product_service/internal/storage"
	"strconv"
)

type Request struct {
	AllProducts   []uint64 `json:"all_products"`
	BrandProducts []uint64 `json:"brand_products"`
	Favorite      []uint64 `json:"favorite"`
	ViewHistory   []uint64 `json:"view_history"`
}

type Response struct {
	AllProducts   []models.Product `json:"all_products"`
	BrandProducts []models.Product `json:"brand_products"`
	Favorite      []models.Product `json:"favorites"`
	ViewHistory   []models.Product `json:"view_history"`
}

type BrandCountRequest struct {
	BrandIDs []uint64 `json:"brand_ids" binding:"required"`
}

type BrandCount struct {
	BrandID uint64 `json:"brand_id"`
	Count   int    `json:"count"`
}

func GetProduct(c *gin.Context) {
	status := c.Param("status")
	count := c.Query("count")
	if count != "" {
		limitInt, err := strconv.Atoi(count)
		if err != nil {
			limitInt = -1
		}
		products, err := storage.GetProducts(status, limitInt)
		if err != nil {
			log.Println("GetProduct: ошибка получения списка товаров", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения списка товаров"})
			return
		}

		c.JSON(200, products)
		return
	}

	var data Request
	var resp Response
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetProduct: ошибка получения данных из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных из запроса"})
		return
	}

	if len(data.AllProducts) != 0 {
		products, err := storage.SelectProduct(data.AllProducts, status)
		if err != nil {
			log.Println("GetProduct: ошибка получения данных о продуктах", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
			return
		}
		resp.AllProducts = products
	}

	if len(data.BrandProducts) != 0 {
		brandProducts, err := storage.SelectProduct(data.BrandProducts, status)
		if err != nil {
			log.Println("GetProduct: ошибка получения данных о продуктах бренда", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
			return
		}
		resp.BrandProducts = brandProducts
	}

	if len(data.Favorite) != 0 {
		productsFavorite, err := storage.SelectProduct(data.Favorite, status)
		if err != nil {
			log.Println("GetProduct: ошибка получения данных об избранных", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об избранных"})
			return
		}

		for i := range productsFavorite {
			productsFavorite[i].IsFavorite = true
		}
		resp.Favorite = productsFavorite
	}

	if len(data.ViewHistory) != 0 {
		productsViewHistory, err := storage.SelectProduct(data.ViewHistory, status)
		if err != nil {
			log.Println("GetProduct: ошибка получения данных об истории просмотра", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об истории просмотра"})
			return
		}
		resp.ViewHistory = productsViewHistory
	}

	c.JSON(200, resp)

}

func GetProductInCategory(c *gin.Context) {
	categoryID := c.Param("id")
	productStatus := c.Param("status")

	products, err := storage.SelectProductsInCategory(categoryID, productStatus)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("GetProductInCategory: ошибка получения данных об товаров из категории", categoryID, err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об истории просмотра"})
		return
	}

	c.JSON(200, products)
}

func GetProductInBrand(c *gin.Context) {
	brandID := c.Param("id")
	productStatus := c.Param("status")

	products, err := storage.SelectProductsInBrand(brandID, productStatus)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("GetProductInBrand: ошибка получения данных об товарах бренда", brandID, err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об истории просмотра"})
		return
	}

	c.JSON(200, products)
}

func CountProductInBrand(c *gin.Context) {
	var req BrandCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("CountProductInBrand: некоректный запрос", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "некорректный запрос"})
		return
	}

	counts, err := storage.GetProductCountsByBrand(req.BrandIDs)
	if err != nil {
		log.Println("CountProductInBrand: ошибка при получении данных", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении данных"})
		return
	}

	var result []BrandCount
	for _, id := range req.BrandIDs {
		result = append(result, BrandCount{
			BrandID: id,
			Count:   counts[id],
		})
	}

	c.JSON(http.StatusOK, result)
}
