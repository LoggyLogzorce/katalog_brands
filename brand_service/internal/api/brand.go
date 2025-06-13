package api

import (
	"brand_service/internal/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type BrandRequest struct {
	BrandsID []uint64 `json:"brands_id"`
}

func GetAllBrands(c *gin.Context) {
	status := c.Param("status")
	count := c.Query("count")

	limitInt, err := strconv.Atoi(count)
	if err != nil {
		limitInt = -1
	}

	brands, err := storage.GetAllBrands(status, limitInt)
	if err != nil {
		log.Println("GetAllBrands: не удалось получить информацию о брендах из бд", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить информацию о брендах"})
		return
	}

	c.JSON(200, brands)
}

func GetBrandInfo(c *gin.Context) {
	var data BrandRequest
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetBrandInfo: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	brand, err := storage.GetBrandInfoById(data.BrandsID)
	if err != nil {
		log.Println("GetBrandInfo: не удалось получить информацию о бренде из бд", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить информацию о бренде"})
		return
	}

	c.JSON(200, brand)
}

func GetBrand(c *gin.Context) {
	brandName := c.Param("name")

	brand, err := storage.GetBrandByName(brandName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(404, gin.H{"error": "бред с данными названием не существует"})
			return
		}
		log.Println("GetBrand: не удалось получить информацию о бренде из бд", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить информацию о бренде"})
		return
	}

	c.JSON(200, brand)
}
