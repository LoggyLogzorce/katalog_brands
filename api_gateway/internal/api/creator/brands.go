package creator

import (
	"api_gateway/internal/api"
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type MyBrand struct {
	ID            uint64  `json:"id"`
	Name          string  `json:"name"`
	LogoURL       string  `json:"logo_url"`
	ProductsCount int     `json:"products_count"`
	AvgRating     float64 `json:"avg_rating"`
	Views         int     `json:"views"`
}

type ProductStatsReq struct {
	ProductIDs []uint64 `json:"product_ids"`
}

func BrandsCreatorHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	url := fmt.Sprintf("/api/v1/brands/creator?creatorID=%s", userID)
	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", url, nil)
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

	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "/api/v1/brands/products", bytes.NewReader(brandsIDJson))
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

	fmt.Println(brandsProductsId)

	productsByBrand := make(map[uint64][]uint64)
	var allProductIDs []uint64
	for _, item := range brandsProductsId {
		productsByBrand[item.BrandID] = item.ProductsId
		allProductIDs = append(allProductIDs, item.ProductsId...)
	}

	buf, _ := json.Marshal(ProductStatsReq{ProductIDs: allProductIDs})
	status, _, body, err = api.ProxyTo(c, "http://localhost:8085", "/api/v1/product_reviews_stats", bytes.NewReader(buf))
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
	status, _, body, err = api.ProxyTo(c, "http://localhost:8082", "/api/v1/product_views_stats", bytes.NewReader(buf))
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
