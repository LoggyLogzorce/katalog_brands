package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"product_service/internal/models"
	"product_service/internal/storage"
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

func GetProduct(c *gin.Context) {
	var data Request
	var resp Response
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetProduct: ошибка получения данных из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных из запроса"})
		return
	}

	if len(data.AllProducts) != 0 {
		products, err := storage.SelectProduct(data.AllProducts)
		if err != nil {
			log.Println("GetProduct: ошибка получения данных о продуктах", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
			return
		}
		resp.AllProducts = products
	}

	if len(data.BrandProducts) != 0 {
		brandProducts, err := storage.SelectProduct(data.BrandProducts)
		if err != nil {
			log.Println("GetProduct: ошибка получения данных о продуктах бренда", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
			return
		}
		resp.BrandProducts = brandProducts
	}

	if len(data.Favorite) != 0 {
		productsFavorite, err := storage.SelectProduct(data.Favorite)
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
		productsViewHistory, err := storage.SelectProduct(data.ViewHistory)
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
	var data Request
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetProduct: ошибка получения данных из запроса", err)
		return
	}

	categoryID := c.Param("id")

	products, err := storage.SelectProductsInCategory(categoryID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("GetProductInCategory: ошибка получения данных об товаров из категории", categoryID, err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об истории просмотра"})
		return
	}

	for i := range products {
		for _, v := range data.Favorite {
			if products[i].ID == v {
				products[i].IsFavorite = true
			}
		}
	}

	c.JSON(200, products)
}
