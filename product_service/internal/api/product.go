package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
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

type BrandRequest struct {
	BrandIDs []uint64 `json:"brand_ids"`
}

type BrandCount struct {
	BrandID uint64 `json:"brand_id"`
	Count   int    `json:"count"`
}

func GetProducts(c *gin.Context) {
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
	var req BrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("CountProductInBrand: некоректный запрос", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "некорректный запрос"})
		return
	}

	role := c.Writer.Header().Get("X-Role")

	counts, err := storage.GetProductCountsByBrand(req.BrandIDs, role)
	if err != nil {
		log.Println("CountProductInBrand: ошибка при получении данных", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка при получении данных"})
		return
	}

	var result []BrandCount
	for _, id := range req.BrandIDs {
		result = append(result, BrandCount{
			BrandID: id,
			Count:   counts[id],
		})
	}

	c.JSON(200, result)
}

func GetProduct(c *gin.Context) {
	productID := c.Param("pId")
	brandID := c.Param("id")
	status := c.Query("status")

	product, err := storage.GetProduct(productID, brandID, status)
	if err != nil {
		log.Println("GetProduct: ошибка при получении данных продукта", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка при получении данных продукта"})
		return
	}

	c.JSON(200, product)
}

func GetProductInBrands(c *gin.Context) {
	var data BrandRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println("GetProductInBrands: некоректный запрос", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "некорректный запрос"})
		return
	}

	products, err := storage.GetProductsInBrands(data.BrandIDs)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("GetProductInBrands: ошибка получения данных об товарах брендов", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об товарах брендов"})
		return
	}

	mp := make(map[uint64][]uint64)

	for _, b := range data.BrandIDs {
		for _, p := range products {
			if b == p.BrandID {
				mp[b] = append(mp[b], p.ID)
			}
		}
	}

	fmt.Println(mp)

	var resp []models.BrandsProductIds

	for k, v := range mp {
		resp = append(resp, models.BrandsProductIds{
			BrandID:    k,
			ProductsId: v,
		})
	}

	c.JSON(200, resp)
}
