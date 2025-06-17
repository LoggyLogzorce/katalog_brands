package api

import (
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func ProductsHandler(c *gin.Context) {
	count := c.Query("count")
	url := fmt.Sprintf("/api/v1/products/approved?count=%s", count)

	status, _, body, err := proxyTo(c, "http://localhost:8083", url, nil)
	if err != nil {
		log.Println("ProductsHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProductsHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var products []models.Product
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("ProductsHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	var brandsID []uint64
	for _, v := range products {
		brandsID = append(brandsID, v.BrandID)
	}

	brandsStruct := models.BrandRequest{
		BrandIDs: brandsID,
	}

	brandsIDJson, err := json.Marshal(brandsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = proxyTo(c, "http://localhost:8084", "/api/v1/brand", bytes.NewReader(brandsIDJson))
	if err != nil {
		log.Println("BrandsHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("BrandsHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var brands []models.Brand
	if err = json.Unmarshal(body, &brands); err != nil {
		log.Println("BrandsHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, body, err = proxyTo(c, "http://localhost:8082", "/api/v1/favorites", nil)
	if err != nil {
		log.Println("ProductsHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProductsHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var favorites []models.Favorite
	if err = json.Unmarshal(body, &favorites); err != nil {
		log.Println("ProductsHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	var productsID []uint64
	for _, v := range products {
		productsID = append(productsID, v.ID)
	}

	productsStruct := models.ReviewsRequest{
		ProductsID: productsID,
	}

	productsIDJson, err := json.Marshal(productsStruct)
	if err != nil {
		log.Println("BrandsHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = proxyTo(c, "http://localhost:8085", "/api/v1/get-reviews", bytes.NewReader(productsIDJson))
	if err != nil {
		log.Println("ProductsHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("ProductsHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var reviews []models.ReviewsResponse
	if err = json.Unmarshal(body, &reviews); err != nil {
		log.Println("ProductsHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	for i := range products {
		for _, v := range brands {
			if products[i].BrandID == v.ID {
				products[i].Brand = v
			}
		}
		for _, v := range favorites {
			if products[i].ID == v.ProductID {
				products[i].IsFavorite = true
			}
		}

		var sum float64
		var cnt int
		for _, rv := range reviews {
			if products[i].ID == rv.ProductID {
				sum += rv.Rating
				cnt++
			}
		}

		products[i].Rating.CountReview = cnt
		if cnt > 0 {
			products[i].Rating.AvgRating = sum / float64(cnt)
		} else {
			products[i].Rating.AvgRating = 0
		}
	}

	c.JSON(http.StatusOK, products)
}

func GetProductHandler(c *gin.Context) {
	brandName := c.Param("name")
	productId := c.Param("id")
	brandUrl := fmt.Sprintf("/api/v1/brand/%s", brandName)
	status, _, body, err := proxyTo(c, "http://localhost:8084", brandUrl, nil)
	if err != nil {
		log.Println("GetProductHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("GetProductHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию о бренде"})
		return
	}

	var brand models.Brand
	if err = json.Unmarshal(body, &brand); err != nil {
		log.Println("GetProductHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	productUrl := fmt.Sprintf("/api/v1/brand/%v/product/%s?status=approved", brand.ID, productId)
	status, _, body, err = proxyTo(c, "http://localhost:8083", productUrl, nil)
	if err != nil {
		log.Println("GetProductHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("GetProductHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию о товаре"})
		return
	}

	var product models.Product
	if err = json.Unmarshal(body, &product); err != nil {
		log.Println("GetProductHandler: ошибка разбора JSON от Product Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Product Service"})
		return
	}

	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	status, _, body, err = proxyTo(c, "http://localhost:8082", "/api/v1/favorites", nil)
	if err != nil {
		log.Println("GetProductHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("GetProductHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
		return
	}

	var favorites []models.Favorite
	if err = json.Unmarshal(body, &favorites); err != nil {
		log.Println("GetProductHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	var reviewsResp []models.ReviewsResponse
	var productID []uint64
	productIDUint, err := strconv.ParseUint(productId, 10, 64)
	if err != nil {
		log.Println("GetProductHandler: ошибка преобразования productId в uint64", err)
	} else {
		productID = append(productID, productIDUint)

		productsStruct := models.ReviewsRequest{
			ProductsID: productID,
		}

		productsIDJson, err := json.Marshal(productsStruct)
		if err != nil {
			log.Println("GetProductHandler: ошибка преобразования brandsID в JSON", err)
		}

		status, _, body, err = proxyTo(c, "http://localhost:8085", "/api/v1/get-reviews", bytes.NewReader(productsIDJson))
		if err != nil {
			log.Println("GetProductHandler: ошибка вызова User Service:", err)
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
			return
		}

		if status != http.StatusOK {
			log.Println("GetProductHandler: User Service вернул статус", status)
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "не удалось получить список категорий"})
			return
		}

		if err = json.Unmarshal(body, &reviewsResp); err != nil {
			log.Println("GetProductHandler: ошибка разбора JSON от User Service:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
			return
		}
	}

	var usersID []uint64
	for _, v := range reviewsResp {
		usersID = append(usersID, v.UserID)
	}

	usersIdData := models.UserDataResponse{
		UsersID: usersID,
	}

	usersIdDataJson, err := json.Marshal(usersIdData)
	if err != nil {
		log.Println("GetProductHandler: ошибка формирования json из userIdData", err)
	}

	role := c.GetString("role")
	c.Request.Header.Set("X-Role", role)

	status, _, body, err = proxyTo(c, "http://localhost:8082", "/api/v1/user-data?count=5", bytes.NewReader(usersIdDataJson))
	if err != nil {
		log.Println("GetProductHandler: ошибка вызова User Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "User Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("GetProductHandler: User Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию для профиля"})
		return
	}

	var users []models.UserData
	if err = json.Unmarshal(body, &users); err != nil {
		log.Println("GetProductHandler: ошибка разбора JSON от User Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа User Service"})
		return
	}

	for _, v := range favorites {
		if product.ID == v.ProductID {
			product.IsFavorite = true
			break
		}
	}

	reviews := make([]models.Review, len(reviewsResp))
	for i, v := range reviewsResp {
		reviews[i].ID = v.ID
		reviews[i].ProductID = v.ProductID
		reviews[i].Rating = v.Rating
		reviews[i].Comment = v.Comment
		reviews[i].CreatedAt = v.CreatedAt

		for _, u := range users {
			if v.UserID == u.UserID {
				reviews[i].User = u
				break
			}
		}
	}

	product.Brand = brand

	resp := models.ProductResponse{
		Product: product,
		Reviews: reviews,
	}

	c.JSON(http.StatusOK, resp)
}
