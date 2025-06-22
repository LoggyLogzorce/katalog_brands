package api

import (
	"brand_service/internal/models"
	"brand_service/internal/storage"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type BrandHandler struct {
	Repo      storage.BrandRepository
	EsService EsService
}

func NewBrandHandler(brandRepo storage.BrandRepository, esService EsService) *BrandHandler {
	return &BrandHandler{Repo: brandRepo, EsService: esService}
}

type BrandRequest struct {
	BrandsID []uint64 `json:"brand_ids"`
}

func (h *BrandHandler) GetAllBrands(c *gin.Context) {
	status := c.Param("status")
	count := c.Query("count")
	creatorID := c.Query("creatorID")

	limitInt, err := strconv.Atoi(count)
	if err != nil {
		limitInt = -1
	}

	brands, err := h.Repo.GetAllBrands(c.Request.Context(), status, creatorID, limitInt)
	if err != nil {
		log.Println("GetAllBrands: не удалось получить информацию о брендах из бд", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить информацию о брендах"})
		return
	}

	c.JSON(200, brands)
}

func (h *BrandHandler) GetBrandInfo(c *gin.Context) {
	var data BrandRequest
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetBrandInfo: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	var brand []models.Brand

	if len(data.BrandsID) == 0 {
		c.JSON(200, brand)
		return
	}

	brand, err := h.Repo.GetBrandInfoById(c.Request.Context(), data.BrandsID)
	if err != nil {
		log.Println("GetBrandInfo: не удалось получить информацию о бренде из бд", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить информацию о бренде"})
		return
	}

	c.JSON(200, brand)
}

func (h *BrandHandler) GetBrand(c *gin.Context) {
	brandName := c.Param("name")
	creatorID := c.Query("creatorID")

	brand, err := h.Repo.GetBrandByName(c.Request.Context(), brandName, creatorID)
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

func (h *BrandHandler) GetBrandByID(c *gin.Context) {
	brandID := c.Param("id")

	brand, err := h.Repo.GetBrandById(c.Request.Context(), brandID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(404, gin.H{"error": "бред с данным названием не существует"})
			return
		}
		log.Println("GetBrand: не удалось получить информацию о бренде из бд", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить информацию о бренде"})
		return
	}

	c.JSON(200, brand)
}

func (h *BrandHandler) CreateBrand(c *gin.Context) {
	var brand models.Brand
	if err := c.ShouldBindBodyWithJSON(&brand); err != nil {
		log.Println("CreateBrand: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	brand.Status = "pending"

	brand, err := h.Repo.CreateBrand(c.Request.Context(), brand)
	if err != nil {
		log.Println("CreateBrand: ошибка создания бренда", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка создания бренда"})
		return
	}

	go h.EsService.CreateUpdateIndex(c.Request.Context(), brand)

	c.JSON(200, gin.H{})
}

func (h *BrandHandler) UpdateBrand(c *gin.Context) {
	var brand models.Brand
	if err := c.ShouldBindBodyWithJSON(&brand); err != nil {
		log.Println("UpdateBrand: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	status := c.Query("status")

	if status != "admin" {
		brand.Status = "pending"
	}

	brand, err := h.Repo.UpdateBrandInfo(c.Request.Context(), brand)
	if err != nil {
		log.Println("UpdateBrand: ошибка обновления данных бренда", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка обновления данных бренда"})
		return
	}

	go h.EsService.CreateUpdateIndex(c.Request.Context(), brand)

	c.JSON(200, gin.H{})
}

func (h *BrandHandler) DeleteBrand(c *gin.Context) {
	brandID := c.Param("id")

	err := h.Repo.DeleteBrand(c.Request.Context(), brandID)
	if err != nil {
		log.Println("DeleteBrand: не удалось удалить бренд", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось удалить бренд"})
		return
	}

	go h.EsService.DeleteIndex(c.Request.Context(), brandID)

	c.JSON(200, gin.H{})
}
