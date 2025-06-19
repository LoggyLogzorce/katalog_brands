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
	"path/filepath"
	"strconv"
	"time"
)

func CreateProductHandler(c *gin.Context) {
	brandName := c.Param("name")

	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
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

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка разбора формы"})
		return
	}

	url := fmt.Sprintf("/api/v1/brand/%s", brandName)
	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", "GET", url, nil)
	if err != nil {
		log.Println("CreateProductHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("CreateProductHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию бренде"})
		return
	}

	var brand models.Brand
	if err = json.Unmarshal(body, &brand); err != nil {
		log.Println("CreateProductHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	files := form.File["images"]
	var urls []models.ProductUrls

	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("brand_%d_%d_%d%s", brand.ID, i, time.Now().Unix(), ext)
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
	product.BrandID = brand.ID
	product.Price = price
	product.CategoryID, err = strconv.ParseUint(categoryID, 10, 64)
	if err != nil {
		log.Println("CreateProductHandler: ошибка конвертации categoryID в uint64", err)
		c.AbortWithStatusJSON(400, gin.H{"error": "ошибка конвертации categoryID в uint64"})
		return
	}
	product.ProductUrls = urls

	productJson, err := json.Marshal(product)
	if err != nil {
		log.Println("CreateProductHandler: ошибка преобразования в json", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "ошибка преобразования в json"})
		return
	}

	url = fmt.Sprintf("/api/v1/product/create")
	status, _, body, err = api.ProxyTo(c, "http://localhost:8083", "POST", url, bytes.NewReader(productJson))
	if err != nil {
		log.Println("CreateProductHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusCreated {
		log.Println("CreateProductHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось создать товар"})
		return
	}

	c.JSON(status, gin.H{})
}

func DeleteProductHandler(c *gin.Context) {
	brandName := c.Param("name")
	productID := c.Param("id")

	url := fmt.Sprintf("/api/v1/brand/%s", brandName)
	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", "GET", url, nil)
	if err != nil {
		log.Println("DeleteProductHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("DeleteProductHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	var brand models.Brand
	if err = json.Unmarshal(body, &brand); err != nil {
		log.Println("DeleteProductHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	url = fmt.Sprintf("/api/v1/brand/%d/product/%s", brand.ID, productID)
	status, _, _, err = api.ProxyTo(c, "http://localhost:8083", "DELETE", url, nil)
	if err != nil {
		log.Println("DeleteProductHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("DeleteProductHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось удалить товар"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateProductHandler(c *gin.Context) {
	brandName := c.Param("name")
	productID := c.Param("id")

	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
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

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка разбора формы"})
		return
	}

	url := fmt.Sprintf("/api/v1/brand/%s", brandName)
	status, _, body, err := api.ProxyTo(c, "http://localhost:8084", "GET", url, nil)
	if err != nil {
		log.Println("UpdateProductHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("UpdateProductHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить информацию бренде"})
		return
	}

	var brand models.Brand
	if err = json.Unmarshal(body, &brand); err != nil {
		log.Println("UpdateProductHandler: ошибка разбора JSON от Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка разбора ответа Brand Service"})
		return
	}

	files := form.File["images"]
	var urls []models.ProductUrls

	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("brand_%d_%d_%d%s", brand.ID, i, time.Now().Unix(), ext)
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
	product.BrandID = brand.ID
	product.Price = price
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

	url = fmt.Sprintf("/api/v1/product/%s", productID)
	status, _, _, err = api.ProxyTo(c, "http://localhost:8083", "PUT", url, bytes.NewReader(productJson))
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
