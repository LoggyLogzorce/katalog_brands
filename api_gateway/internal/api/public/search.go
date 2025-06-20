package public

import (
	"api_gateway/internal/api"
	"api_gateway/internal/es"
	"api_gateway/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	esClient, _ = es.NewClient()
	indexer     = es.NewIndexer(esClient)
)

func SearchAllHandler(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр q обязателен"})
		return
	}
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	from, _ := strconv.Atoi(c.DefaultQuery("from", "0"))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var (
		productsEs   []es.ProductDoc
		brandsEs     []es.BrandDoc
		categoriesEs []es.CategoryDoc
		prodErr      error
		brandErr     error
		catErr       error
	)

	wg.Add(3)
	go func() {
		defer wg.Done()
		productsEs, prodErr = indexer.SearchProducts(ctx, q, size, from)
	}()
	go func() {
		defer wg.Done()
		brandsEs, brandErr = indexer.SearchBrands(ctx, q, size, from)
	}()
	go func() {
		defer wg.Done()
		categoriesEs, catErr = indexer.SearchCategories(ctx, q, size, from)
	}()
	wg.Wait()

	if prodErr != nil || brandErr != nil || catErr != nil {
		log.Println("SearchAll error:", prodErr, brandErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка поиска"})
		return
	}

	products, status, err := productsDto(c, productsEs)
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{"error": err})
		return
	}

	brands, status, err := brandsDto(c, brandsEs)
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{"error": err})
		return
	}

	categories, status, err := categoriesDto(c, categoriesEs)
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products":   gin.H{"total": len(products), "items": products},
		"brands":     gin.H{"total": len(brands), "items": brands},
		"categories": gin.H{"total": len(categories), "items": categories},
	})
}

func productsDto(c *gin.Context, productsEs []es.ProductDoc) ([]models.Product, int, error) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, body, err := api.ProxyTo(c, "http://localhost:8082", "", "/api/v1/favorites", nil)
	if err != nil {
		log.Println("FavoriteHandler: ошибка вызова User Service:", err)
		return nil, http.StatusBadGateway, err
	}

	if status != http.StatusOK {
		log.Println("FavoriteHandler: User Service вернул статус", status)
		return nil, status, errors.New(fmt.Sprintf("user Service вернул статус %d", status))
	}

	var favorites []models.Favorite
	if err = json.Unmarshal(body, &favorites); err != nil {
		log.Println("FavoriteHandler: ошибка разбора JSON от User Service:", err)
		return nil, http.StatusInternalServerError, err
	}

	var productsID []uint64
	for _, v := range productsEs {
		if v.Status == "approved" {
			id, err := strconv.ParseUint(v.ID, 10, 64)
			if err != nil {
				log.Println("SearchAllHandler: ошибка преобразования productID в uint64", err)
				continue
			}
			productsID = append(productsID, id)
		}
	}

	productIdJson, err := json.Marshal(models.ProductRequest{AllProducts: productsID})
	if err != nil {
		log.Println("SearchAllHandler: ошибка маршалинга", err)
		return nil, http.StatusInternalServerError, err
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "", "/api/v1/products/approved", bytes.NewBuffer(productIdJson))
	if err != nil {
		log.Println("FavoriteHandler: ошибка вызова Product Service:", err)
		return nil, http.StatusBadGateway, err
	}

	if status != http.StatusOK {
		log.Println("FavoriteHandler: Product Service вернул статус", status)
		return nil, status, errors.New(fmt.Sprintf("user Service вернул статус %d", status))
	}

	var products models.ProfileProductResponse
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("FavoriteHandler: ошибка разбора JSON от Product Service:", err)
		return nil, http.StatusInternalServerError, err
	}

	productsID = []uint64{}
	var brandsID []uint64

	for _, v := range products.AllProducts {
		productsID = append(productsID, v.ID)
		brandsID = append(brandsID, v.BrandID)
	}

	productsIDJson, err := json.Marshal(models.ReviewsRequest{ProductsID: productsID})
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8085", "", "/api/v1/get-reviews", bytes.NewReader(productsIDJson))
	if err != nil {
		log.Println("ProductsHandler: ошибка вызова User Service:", err)
		return nil, http.StatusBadGateway, err
	}

	if status != http.StatusOK {
		log.Println("ProductsHandler: User Service вернул статус", status)
		return nil, status, errors.New(fmt.Sprintf("user Service вернул статус %d", status))
	}

	var reviews []models.ReviewsResponse
	if err = json.Unmarshal(body, &reviews); err != nil {
		log.Println("ProductsHandler: ошибка разбора JSON от User Service:", err)
		return nil, http.StatusInternalServerError, err
	}

	brandsStruct := models.BrandRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(brandsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8084", "", "/api/v1/brand", bytes.NewReader(brandsIDJson))
	if err != nil {
		log.Println("BrandsHandler: ошибка вызова Brand Service:", err)
		return nil, http.StatusBadGateway, err
	}

	if status != http.StatusOK {
		log.Println("BrandsHandler: Brand Service вернул статус", status)
		return nil, status, errors.New(fmt.Sprintf("user Service вернул статус %d", status))
	}

	var brands []models.Brand
	if err = json.Unmarshal(body, &brands); err != nil {
		log.Println("BrandsHandler: ошибка разбора JSON от Brand Service:", err)
		return nil, http.StatusInternalServerError, err
	}

	for i, v := range products.AllProducts {
		for _, b := range brands {
			if v.BrandID == b.ID {
				products.AllProducts[i].Brand = b
			}
		}
		for _, f := range favorites {
			if products.AllProducts[i].ID == f.ProductID {
				products.AllProducts[i].IsFavorite = true
			}
		}

		var sum float64
		var cnt int
		for _, rv := range reviews {
			if products.AllProducts[i].ID == rv.ProductID {
				sum += rv.Rating
				cnt++
			}
		}

		products.AllProducts[i].Rating.CountReview = cnt
		if cnt > 0 {
			products.AllProducts[i].Rating.AvgRating = sum / float64(cnt)
		} else {
			products.AllProducts[i].Rating.AvgRating = 0
		}
	}

	return products.AllProducts, http.StatusOK, nil
}

func brandsDto(c *gin.Context, brandsEs []es.BrandDoc) ([]models.Brand, int, error) {
	var brandsID []uint64
	for _, v := range brandsEs {
		if v.Status == "approved" {
			brandIdUint, err := strconv.ParseUint(v.ID, 10, 64)
			if err != nil {
				log.Println("Brands: ошибка конвертации brandID в uint64", err)
				continue
			}
			brandsID = append(brandsID, brandIdUint)
		}
	}

	brandsStruct := models.BrandRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(brandsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", "", "/api/v1/brand", bytes.NewReader(brandsIDJson))
	if err != nil {
		log.Println("BrandsHandler: ошибка вызова Brand Service:", err)
		return nil, http.StatusBadGateway, err
	}

	if status != http.StatusOK {
		log.Println("BrandsHandler: Brand Service вернул статус", status)
		return nil, status, errors.New(fmt.Sprintf("user Service вернул статус %d", status))
	}

	var brands []models.Brand
	if err = json.Unmarshal(body, &brands); err != nil {
		log.Println("BrandsHandler: ошибка разбора JSON от Brand Service:", err)
		return nil, http.StatusInternalServerError, err
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "", "/api/v1/brands/count-product", bytes.NewReader(brandsIDJson))
	if err != nil {
		log.Println("BrandsHandler: ошибка вызова Product Service:", err)
		return nil, http.StatusBadGateway, err
	}

	if status != http.StatusOK {
		log.Println("BrandsHandler: Product Service вернул статус", status)
		return nil, status, errors.New(fmt.Sprintf("user Service вернул статус %d", status))
	}

	var brandsCount []models.BrandCount
	if err = json.Unmarshal(body, &brandsCount); err != nil {
		log.Println("BrandsHandler: ошибка разбора JSON от Product Service:", err)
		return nil, http.StatusInternalServerError, err
	}

	for i, v := range brands {
		for _, e := range brandsCount {
			if v.ID == e.BrandID {
				brands[i].ProductCount = e.Count
			}
		}
	}

	return brands, http.StatusOK, nil
}

func categoriesDto(c *gin.Context, categoriesEs []es.CategoryDoc) ([]models.Category, int, error) {
	url := fmt.Sprintf("/api/v1/categories")
	status, _, body, err := api.ProxyTo(c, "http://localhost:8083", "GET", url, nil)
	if err != nil {
		log.Println("CategoryHandler: ошибка вызова Product Service:", err)
		return nil, http.StatusBadGateway, err
	}

	if status != http.StatusOK {
		log.Println("CategoryHandler: Product Service вернул статус", status)
		return nil, http.StatusBadGateway, err
	}

	var categories []models.Category
	if err = json.Unmarshal(body, &categories); err != nil {
		log.Println("CategoryHandler: ошибка разбора JSON от Product Service:", err)
		return nil, http.StatusInternalServerError, err
	}

	var categoriesNew []models.Category
	for _, v := range categories {
		for _, n := range categoriesEs {
			if n.ID == strconv.FormatUint(v.ID, 10) {
				categoriesNew = append(categoriesNew, v)
			}
		}
	}

	return categoriesNew, http.StatusOK, nil
}
