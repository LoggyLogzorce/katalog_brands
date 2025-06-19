package admin

import (
	"api_gateway/internal/api"
	"api_gateway/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func GetProductsAdminHandler(c *gin.Context) {
	url := fmt.Sprintf("/api/v1/products/admin?count=all")

	status, _, body, err := api.ProxyTo(c, "http://localhost:8083", "", url, nil)
	if err != nil {
		log.Println("GetProductsAdminHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("GetProductsAdminHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var products []models.Product
	if err = json.Unmarshal(body, &products); err != nil {
		log.Println("GetProductsAdminHandler: ошибка разбора JSON от Product Service:", err)
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
		log.Println("GetProductsAdminHandler: ошибка преобразования brandsID в JSON", err)
	}

	status, _, body, err = api.ProxyTo(c, "http://localhost:8084", "", "/api/v1/brand", bytes.NewReader(brandsIDJson))
	if err != nil {
		log.Println("GetProductsAdminHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("GetProductsAdminHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var brands []models.Brand
	if err = json.Unmarshal(body, &brands); err != nil {
		log.Println("GetProductsAdminHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	for i := range products {
		for _, v := range brands {
			if products[i].BrandID == v.ID {
				products[i].Brand = v
			}
		}
	}

	c.JSON(http.StatusOK, products)
}

func CreateProductAdminHandler(c *gin.Context) {
	name := c.PostForm("name")
	brandIdStr := c.PostForm("brand_id")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	categoryID := c.PostForm("category_id")
	prStatus := c.PostForm("status")

	if name == "" || description == "" || priceStr == "" || categoryID == "" || brandIdStr == "" || prStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "все поля обязательны"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректная цена"})
		return
	}

	brandID, err := strconv.ParseUint(brandIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id бренда"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка разбора формы"})
		return
	}

	files := form.File["images"]
	var urls []models.ProductUrls

	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("brand_%d_%d_%d%s", brandID, i, time.Now().Unix(), ext)
		dst := filepath.Join("web", "static", "img", "products", newFilename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить логотип"})
			return
		}
		urls = append(urls, models.ProductUrls{
			ProductID: 0,
			Url:       "img/products/" + newFilename,
		})
	}

	var product models.Product
	product.Name = name
	product.Description = description
	product.BrandID = brandID
	product.Price = price
	product.Status = prStatus
	product.CategoryID, err = strconv.ParseUint(categoryID, 10, 64)
	if err != nil {
		log.Println("CreateProductAdminHandler: ошибка конвертации categoryID в uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка конвертации categoryID в uint64"})
		return
	}
	product.ProductUrls = urls

	productJson, err := json.Marshal(product)
	if err != nil {
		log.Println("CreateProductAdminHandler: ошибка преобразования в json", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка преобразования в json"})
		return
	}

	url := fmt.Sprintf("/api/v1/product/create")
	status, _, _, err := api.ProxyTo(c, "http://localhost:8083", "POST", url, bytes.NewReader(productJson))
	if err != nil {
		log.Println("CreateProductAdminHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusCreated {
		log.Println("CreateProductAdminHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось создать товар"})
		return
	}

	c.JSON(status, gin.H{})
}

func UpdateProductAdminHandler(c *gin.Context) {
	productID := c.Param("id")
	name := c.PostForm("name")
	brandIdStr := c.PostForm("brand_id")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	strStatus := c.PostForm("status")
	categoryID := c.PostForm("category_id")

	if name == "" || description == "" || priceStr == "" || categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "все поля обязательны"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректная цена"})
		return
	}

	brandID, err := strconv.ParseUint(brandIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id бренда"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка разбора формы"})
		return
	}

	files := form.File["images"]
	var urls []models.ProductUrls

	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("brand_%d_%d_%d%s", brandID, i, time.Now().Unix(), ext)
		dst := filepath.Join("web", "static", "img", "products", newFilename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить логотип"})
			return
		}
		urls = append(urls, models.ProductUrls{
			ProductID: 0,
			Url:       "img/products/" + newFilename,
		})
	}

	var product models.Product
	product.Name = name
	product.Description = description
	product.BrandID = brandID
	product.Price = price
	product.Status = strStatus
	product.CategoryID, err = strconv.ParseUint(categoryID, 10, 64)
	if err != nil {
		log.Println("UpdateProductHandler: ошибка конвертации categoryID в uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка конвертации categoryID в uint64"})
		return
	}
	product.ProductUrls = urls

	productJson, err := json.Marshal(product)
	if err != nil {
		log.Println("UpdateProductHandler: ошибка преобразования в json", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка преобразования в json"})
		return
	}

	url := fmt.Sprintf("/api/v1/product/%s", productID)
	status, _, _, err := api.ProxyTo(c, "http://localhost:8083", "PUT", url, bytes.NewReader(productJson))
	if err != nil {
		log.Println("UpdateProductHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("UpdateProductHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось обновить информации о товаре"})
		return
	}

	c.JSON(status, gin.H{})
}

func DeleteProductAdminHandler(c *gin.Context) {
	productID := c.Param("id")

	url := fmt.Sprintf("/api/v1/brand/%s/product/%s?status=admin", "nil", productID)
	status, _, _, err := api.ProxyTo(c, "http://localhost:8083", "DELETE", url, nil)
	if err != nil {
		log.Println("DeleteProductAdminHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("DeleteProductAdminHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось удалить товар"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
