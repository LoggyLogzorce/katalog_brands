package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"product_service/internal/es"
	"product_service/internal/models"
	"product_service/internal/storage"
	"strconv"
	"time"
)

type ProductHandler struct {
	PrRepo storage.ProductRepository
	EsRepo es.IndexerRepository
}

func NewProductHandler(pr storage.ProductRepository, es es.IndexerRepository) *ProductHandler {
	return &ProductHandler{PrRepo: pr, EsRepo: es}
}

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

func (h *ProductHandler) GetProductsHandler(c *gin.Context) {
	status := c.Param("status")
	count := c.Query("count")
	if count != "" {
		limitInt, err := strconv.Atoi(count)
		if err != nil {
			limitInt = -1
		}
		products, err := h.PrRepo.GetProducts(c.Request.Context(), status, limitInt)
		if err != nil {
			log.Println("GetProductHandler: ошибка получения списка товаров", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения списка товаров"})
			return
		}

		c.JSON(200, products)
		return
	}

	var data Request
	var resp Response
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("GetProductHandler: ошибка получения данных из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных из запроса"})
		return
	}

	if len(data.AllProducts) != 0 {
		products, err := h.PrRepo.SelectProduct(c.Request.Context(), data.AllProducts, status)
		if err != nil {
			log.Println("GetProductHandler: ошибка получения данных о продуктах", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
			return
		}
		resp.AllProducts = products
	}

	if len(data.BrandProducts) != 0 {
		brandProducts, err := h.PrRepo.SelectProduct(c.Request.Context(), data.BrandProducts, status)
		if err != nil {
			log.Println("GetProductHandler: ошибка получения данных о продуктах бренда", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных о продуктах"})
			return
		}
		resp.BrandProducts = brandProducts
	}

	if len(data.Favorite) != 0 {
		productsFavorite, err := h.PrRepo.SelectProduct(c.Request.Context(), data.Favorite, status)
		if err != nil {
			log.Println("GetProductHandler: ошибка получения данных об избранных", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об избранных"})
			return
		}

		for i := range productsFavorite {
			productsFavorite[i].IsFavorite = true
		}
		resp.Favorite = productsFavorite
	}

	if len(data.ViewHistory) != 0 {
		productsViewHistory, err := h.PrRepo.SelectProduct(c.Request.Context(), data.ViewHistory, status)
		if err != nil {
			log.Println("GetProductHandler: ошибка получения данных об истории просмотра", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об истории просмотра"})
			return
		}
		resp.ViewHistory = productsViewHistory
	}

	c.JSON(200, resp)

}

func (h *ProductHandler) GetProductInCategoryHandler(c *gin.Context) {
	categoryID := c.Param("id")
	productStatus := c.Param("status")

	products, err := h.PrRepo.SelectProductsInCategory(c.Request.Context(), categoryID, productStatus)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("GetProductInCategoryHandler: ошибка получения данных об товаров из категории", categoryID, err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об истории просмотра"})
		return
	}

	c.JSON(200, products)
}

func (h *ProductHandler) GetProductInBrandHandler(c *gin.Context) {
	brandID := c.Param("id")
	productStatus := c.Param("status")

	products, err := h.PrRepo.SelectProductsInBrand(c.Request.Context(), brandID, productStatus)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("GetProductInBrandHandler: ошибка получения данных об товарах бренда", brandID, err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка получения данных об истории просмотра"})
		return
	}

	c.JSON(200, products)
}

func (h *ProductHandler) CountProductInBrandHandler(c *gin.Context) {
	var req BrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("CountProductInBrandHandler: некоректный запрос", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "некорректный запрос"})
		return
	}

	role := c.Writer.Header().Get("X-Role")

	counts, err := h.PrRepo.GetProductCountsByBrand(c.Request.Context(), req.BrandIDs, role)
	if err != nil {
		log.Println("CountProductInBrandHandler: ошибка при получении данных", err)
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

func (h *ProductHandler) GetProductHandler(c *gin.Context) {
	productID := c.Param("pId")
	brandID := c.Param("id")
	status := c.Query("status")

	product, err := h.PrRepo.GetProduct(c.Request.Context(), productID, brandID, status)
	if err != nil {
		log.Println("GetProductHandler: ошибка при получении данных продукта", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка при получении данных продукта"})
		return
	}

	c.JSON(200, product)
}

func (h *ProductHandler) GetProductInBrandsHandler(c *gin.Context) {
	var data BrandRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println("GetProductInBrandsHandler: некоректный запрос", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "некорректный запрос"})
		return
	}

	products, err := h.PrRepo.GetProductsInBrands(c.Request.Context(), data.BrandIDs)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("GetProductInBrandsHandler: ошибка получения данных об товарах брендов", err)
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

	var resp []models.BrandsProductIds

	for k, v := range mp {
		resp = append(resp, models.BrandsProductIds{
			BrandID:    k,
			ProductsId: v,
		})
	}

	c.JSON(200, resp)
}

func (h *ProductHandler) CreateProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Println("CreateProductHandler:", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "не удалось получить данные"})
		return
	}

	if product.Status == "" {
		product.Status = "pending"
	}
	product.CreatedAt = time.Now()

	// сохраняем в БД
	product, err := h.PrRepo.CreateProduct(c.Request.Context(), product, product.ProductUrls)
	if err != nil {
		log.Println("CreateProductHandler:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить товар"})
		return
	}

	go func(data models.Product) {
		p, err := h.PrRepo.GetProduct(c.Request.Context(), strconv.FormatUint(data.ID, 10), strconv.FormatUint(data.BrandID, 10), data.Status)
		if err != nil {
			log.Println("UpdateProductHandler: не удалось получить товар для индексации", err)
			return
		}

		doc := es.ProductDoc{
			ID:          fmt.Sprint(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category.Name,
			Price:       p.Price,
			Photo:       p.ProductUrls[0].Url,
			Status:      p.Status,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.EsRepo.IndexProduct(ctx, doc); err != nil {
			log.Println("CreateProductHandler: ES indexing failed:", err)
			return
		}
		log.Println("CreateProductHandler: ES indexing successfully")
	}(product)

	c.JSON(http.StatusCreated, gin.H{})
}

func (h *ProductHandler) UpdateProductHandler(c *gin.Context) {
	var data models.Product
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		log.Println("UpdateProductHandler: не удалось получить данные из запроса", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "не удалось получить данные из запроса"})
		return
	}

	productID := c.Param("id")
	productIdUint, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		log.Println("UpdateProductHandler: ошибка преобразования productID в uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка преобразования productID в uint64"})
		return
	}

	if data.Status == "" {
		data.Status = "pending"
	}

	data.ID = productIdUint

	for i := range data.ProductUrls {
		data.ProductUrls[i].ProductID = data.ID
	}

	err = h.PrRepo.UpdateProduct(c.Request.Context(), data)
	if err != nil {
		log.Println("UpdateProductHandler: не удалось обновить данные товара", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось обновить данные товара"})
		return
	}

	go func() {
		p, err := h.PrRepo.GetProduct(c.Request.Context(), strconv.FormatUint(data.ID, 10), strconv.FormatUint(data.BrandID, 10), data.Status)
		if err != nil {
			log.Println("UpdateProductHandler: не удалось получить товар для индексации", err)
			return
		}

		doc := es.ProductDoc{
			ID:          fmt.Sprint(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category.Name,
			Price:       p.Price,
			Photo:       p.ProductUrls[0].Url,
			Status:      p.Status,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.EsRepo.IndexProduct(ctx, doc); err != nil {
			log.Println("CreateProductHandler: ES indexing failed:", err)
			return
		}
		log.Println("CreateProductHandler: ES indexing successfully")
	}()

	c.JSON(200, gin.H{})
}

func (h *ProductHandler) DeleteProductHandler(c *gin.Context) {
	brandID := c.Param("id")
	productID := c.Param("pId")
	status := c.Query("status")

	err := h.PrRepo.Delete(c.Request.Context(), brandID, productID, status)
	if err != nil {
		log.Println("DeleteProductHandler: не удалось удалить товар", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось удалить товар"})
		return
	}

	go func(id string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.EsRepo.DeleteProduct(ctx, id); err != nil {
			log.Println("DeleteProductHandler: ошибка удаления из ES: doc", err)
			return
		}
		log.Println("DeleteProductHandler: успешное удаление из ES:", id)
	}(productID)

	c.JSON(200, gin.H{})
}
