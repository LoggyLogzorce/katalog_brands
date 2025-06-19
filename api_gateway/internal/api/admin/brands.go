package admin

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

func BrandsAdminHandler(c *gin.Context) {
	userID := c.GetString("userID")
	c.Request.Header.Set("X-User-ID", userID)

	url := fmt.Sprintf("/api/v1/brands/all")
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

	c.JSON(200, brands)
}

func DeleteBrandHandler(c *gin.Context) {
	brandID := c.Param("id")
	url := fmt.Sprintf("/api/v1/brand/%s", brandID)
	status, _, _, err := api.ProxyTo(c, "http://localhost:8084", "", url, nil)
	if err != nil {
		log.Println("DeleteBrandHandler: ошибка вызова Brand Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Brand Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("DeleteBrandHandler: Brand Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось получить список брендов"})
		return
	}

	c.JSON(200, gin.H{})
}

func CreateBrandAdminHandler(c *gin.Context) {
	brandName := c.PostForm("name")
	brandDesc := c.PostForm("description")
	brandCreator := c.PostForm("creator")
	brandStatus := c.PostForm("status")
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
	brandCreatorUint, err := strconv.ParseUint(brandCreator, 10, 64)
	if err != nil {
		log.Println("CreateBrandAdminHandler: не удалось преобразовать brandCreator в uint64", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "не удалось преобразовать brandCreator в uint64"})
		return
	}
	brand.CreatorID = brandCreatorUint
	brand.Status = brandStatus

	brandJson, err := json.Marshal(brand)
	url = fmt.Sprintf("/api/v1/brand/create")
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

func UpdateBrandAdminHandler(c *gin.Context) {
	brandID := c.Param("id")
	newName := c.PostForm("name")
	newDesc := c.PostForm("description")
	newCreator := c.PostForm("creator")
	newStatus := c.PostForm("status")
	var fileError error

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		fileError = err
		log.Println("UpdateBrandHandler: не удалось считать логотип из формы", err)
	}

	url := fmt.Sprintf("/api/v1/brand/get/%s", brandID)
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
	brand.Status = newStatus
	newCreatorUint, err := strconv.ParseUint(newCreator, 10, 64)
	if err != nil {
		log.Println("UpdateBrandAdminHandler: не удалось преобразовать newCreator в uint64", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "не удалось преобразовать newCreator в uint64"})
		return
	}
	brand.CreatorID = newCreatorUint

	brandJson, err := json.Marshal(brand)
	url = fmt.Sprintf("/api/v1/brand/%s?status=admin", brand.Name)
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
