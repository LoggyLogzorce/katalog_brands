package creator

import (
	"api_gateway/internal/api"
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type MyBrand struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	LogoURL       string    `json:"logo_url"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	ProductsCount int       `json:"products_count"`
	AvgRating     float64   `json:"avg_rating"`
	Views         int       `json:"views"`
}

type ProductStatsReq struct {
	ProductIDs []uint64 `json:"product_ids"`
}

func BrandsCreatorHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	url := fmt.Sprintf("/api/v1/brands/creator?creatorID=%s", userID)
	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", "GET", url, nil)
	if err != nil {
		log.Println("BrandsCreatorHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandsCreatorHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var brands []models.Brand
	if err = json.Unmarshal(body, &brands); err != nil {
		log.Println("BrandsCreatorHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	var brandsID []uint64
	for _, v := range brands {
		brandsID = append(brandsID, v.ID)
	}

	productsCount := models.BrandRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(productsCount)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "GET", "/api/v1/brands/products", bytes.NewReader(brandsIDJson))
	if err != nil {
		log.Println("BrandsCreatorHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandsCreatorHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить количество товаров брендов"})
		return
	}

	var brandsProductsId []models.BrandsProductIds
	if err = json.Unmarshal(body, &brandsProductsId); err != nil {
		log.Println("BrandsCreatorHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	productsByBrand := make(map[uint64][]uint64)
	var allProductIDs []uint64
	for _, item := range brandsProductsId {
		productsByBrand[item.BrandID] = item.ProductsId
		allProductIDs = append(allProductIDs, item.ProductsId...)
	}

	buf, _ := json.Marshal(ProductStatsReq{ProductIDs: allProductIDs})
	status, _, body, err = api.ProxyTo(c, "http://localhost:8085", "GET", "/api/v1/product_reviews_stats", bytes.NewReader(buf))
	if err != nil {
		log.Println("BrandsCreatorHandler: ошибка вызова Review Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Review Service недоступен"})
		return
	}
	if status != http.StatusOK {
		log.Println("BrandsCreatorHandler: Review Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить количество товаров брендов"})
		return
	}

	var revStats []models.ProductReviewStat
	if err = json.Unmarshal(body, &revStats); err != nil {
		log.Println("BrandsCreatorHandler: ошибка разбора JSON от Review Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Review Service"})
		return
	}

	rsMap := make(map[uint64]models.ProductReviewStat)
	for _, s := range revStats {
		rsMap[s.ProductID] = s
	}

	buf, _ = json.Marshal(ProductStatsReq{ProductIDs: allProductIDs})
	status, _, body, err = api.ProxyTo(c, "http://localhost:8082", "GET", "/api/v1/product_views_stats", bytes.NewReader(buf))
	if err != nil {
		log.Println("BrandsCreatorHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}
	if status != http.StatusOK {
		log.Println("BrandsCreatorHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить количество товаров брендов"})
		return
	}

	var vwStats []models.ProductViewStat
	if err = json.Unmarshal(body, &vwStats); err != nil {
		log.Println("BrandsCreatorHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	vsMap := make(map[uint64]int)
	for _, s := range vwStats {
		vsMap[s.ProductID] = s.Views
	}

	var result []MyBrand
	for _, b := range brands {
		pIDs := productsByBrand[b.ID]
		var sumRating float64
		var cntReview, sumViews int
		for _, pid := range pIDs {
			if st, ok := rsMap[pid]; ok {
				sumRating += st.AvgRating * float64(st.CountReview)
				cntReview += st.CountReview
			}
			sumViews += vsMap[pid]
		}
		avg := 0.0
		if cntReview > 0 {
			avg = sumRating / float64(cntReview)
		}
		result = append(result, MyBrand{
			ID:            b.ID,
			Name:          b.Name,
			LogoURL:       b.LogoUrl,
			ProductsCount: len(pIDs),
			AvgRating:     avg,
			Views:         sumViews,
		})
	}

	c.JSON(http.StatusOK, gin.H{"brands": result})
}

func BrandHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	brandName := c.Param("name")

	url := fmt.Sprintf("/api/v1/brand/%s?creatorID=%s", brandName, userID)
	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", "GET", url, nil)
	if err != nil {
		log.Println("BrandsCreatorHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandsCreatorHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var brand models.Brand
	if err = json.Unmarshal(body, &brand); err != nil {
		log.Println("BrandsCreatorHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	url = fmt.Sprintf("/api/v1/brand/%v/products/all", brand.ID)
	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "GET", url, nil)
	if err != nil {
		log.Println("BrandHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию бренде"})
		return
	}

	var products []models.Product
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("BrandHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	var allProductIDs []uint64
	for _, v := range products {
		allProductIDs = append(allProductIDs, v.ID)
	}

	buf, _ := json.Marshal(ProductStatsReq{ProductIDs: allProductIDs})
	status, _, body, err = api.ProxyTo(c, "http://localhost:8082", "GET", "/api/v1/product_views_stats", bytes.NewReader(buf))
	if err != nil {
		log.Println("BrandsCreatorHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}
	if status != http.StatusOK {
		log.Println("BrandsCreatorHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить количество товаров брендов"})
		return
	}

	var vwStats []models.ProductViewStat
	if err = json.Unmarshal(body, &vwStats); err != nil {
		log.Println("BrandsCreatorHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	buf, _ = json.Marshal(ProductStatsReq{ProductIDs: allProductIDs})
	status, _, body, err = api.ProxyTo(c, "http://localhost:8085", "GET", "/api/v1/product_reviews_stats", bytes.NewReader(buf))
	if err != nil {
		log.Println("BrandsCreatorHandler: ошибка вызова Review Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Review Service недоступен"})
		return
	}
	if status != http.StatusOK {
		log.Println("BrandsCreatorHandler: Review Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить количество товаров брендов"})
		return
	}

	var revStats []models.ProductReviewStat
	if err = json.Unmarshal(body, &revStats); err != nil {
		log.Println("BrandsCreatorHandler: ошибка разбора JSON от Review Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Review Service"})
		return
	}

	var resultBrand MyBrand
	resultBrand.ID = brand.ID
	resultBrand.Name = brand.Name
	resultBrand.Description = brand.Description
	resultBrand.LogoURL = brand.LogoUrl
	resultBrand.Status = brand.Status
	resultBrand.CreatedAt = brand.CreatedAt
	resultBrand.ProductsCount = len(allProductIDs)

	for _, v := range vwStats {
		resultBrand.Views += v.Views
	}
	var avgRating float64
	var revCount float64
	for _, v := range revStats {
		avgRating += v.AvgRating * float64(v.CountReview)
		revCount += float64(v.CountReview)
	}
	if revCount == 0 {
		resultBrand.AvgRating = 0
	} else {
		resultBrand.AvgRating = avgRating / revCount
	}

	result := gin.H{
		"brand":   resultBrand,
		"product": products,
	}

	c.JSON(http.StatusOK, result)
}

func UpdateBrandHandler(c *gin.Context) {
	brandName := c.Param("name")
	newName := c.PostForm("name")
	newDesc := c.PostForm("description")
	var fileError error

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		fileError = err
		log.Println("UpdateBrandHandler: не удалось считать логотип из формы", err)
	}

	url := fmt.Sprintf("/api/v1/brand/%s", brandName)
	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", "GET", url, nil)
	if err != nil {
		log.Println("UpdateBrandHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("UpdateBrandHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию бренде"})
		return
	}

	var brand models.Brand
	if err = json.Unmarshal(body, &brand); err != nil {
		log.Println("UpdateBrandHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	if !errors.Is(fileError, http.ErrMissingFile) {
		ext := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("brand_%d_%d%s", brand.ID, time.Now().Unix(), ext)
		dst := filepath.Join("web", "static", "img", "logo", newFilename)

		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить логотип"})
			return
		}
		fmt.Println(dst)
		brand.LogoUrl = "img/logo/" + newFilename
	}

	brand.Name = newName
	brand.Description = newDesc

	brandJson, err := json.Marshal(brand)

	status, _, _, err = api.ProxyTo(c, "http://localhost:8084", "PUT", url, bytes.NewReader(brandJson))
	if err != nil {
		log.Println("UpdateBrandHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("UpdateBrandHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось обновить информацию о бренде"})
		return
	}

	c.JSON(http.StatusOK, brand)
}

func CreateBrandHandler(c *gin.Context) {
	brandName := c.PostForm("name")
	brandDesc := c.PostForm("description")
	var fileError error

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		fileError = err
		log.Println("CreateBrandHandler: не удалось считать логотип из формы", err)
	}

	url := fmt.Sprintf("/api/v1/brand/%s", brandName)
	status, _, _, err := api.ProxyTo(c, "http://localhost:8084", "GET", url, nil)
	if err != nil {
		log.Println("CreateBrandHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusNotFound {
		if status == http.StatusOK {
			log.Println("CreateBrandHandler: Brand Service вернул статус", status)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Бренд с таким названием уже сущесвует"})
			return
		}
		log.Println("CreateBrandHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию бренде"})
		return
	}

	var brand models.Brand
	if !errors.Is(fileError, http.ErrMissingFile) {
		ext := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("brand_%d_%d%s", brand.ID, time.Now().Unix(), ext)
		dst := filepath.Join("web", "static", "img", "logo", newFilename)

		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить логотип"})
			return
		}
		fmt.Println(dst)
		brand.LogoUrl = "img/logo/" + newFilename
	}

	userID := c.GetString("userID")

	userIdUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		log.Println("CreateBrandHandler: ошибка преобразования userID в uint64", err)
	} else {
		brand.CreatorID = userIdUint
	}

	brand.Name = brandName
	brand.Description = brandDesc

	brandJson, err := json.Marshal(brand)
	status, _, _, err = api.ProxyTo(c, "http://localhost:8084", "POST", url, bytes.NewReader(brandJson))
	if err != nil {
		log.Println("CreateBrandHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("CreateBrandHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось обновить информацию о бренде"})
		return
	}

	c.JSON(http.StatusOK, brand)
}
