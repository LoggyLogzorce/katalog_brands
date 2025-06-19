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

func CreateCategoryAdminHandler(c *gin.Context) {
	categoryName := c.PostForm("category_name")

	var category models.Category
	category.Name = categoryName

	fileHeader, err := c.FormFile("logoImage")
	if err == nil {
		ext := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("category_%s_%d%s", category.Name, time.Now().Unix(), ext)
		dst := filepath.Join("web", "static", "img", "category", newFilename)
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить изображение"})
			return
		}
		category.Photo = "img/category/" + newFilename
	} else if !errors.Is(err, http.ErrMissingFile) {
		log.Println("CreateCategoryAdminHandler: ошибка чтения logoImage", err)
	}

	categoryJson, err := json.Marshal(category)
	url := fmt.Sprintf("/api/v1/category/create")
	status, _, _, err := api.ProxyTo(c, "http://localhost:8083", "", url, bytes.NewReader(categoryJson))
	if err != nil {
		log.Println("CreateCategoryAdminHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusCreated {
		log.Println("CreateCategoryAdminHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось создать категорию"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func UpdateCategoryAdminHandler(c *gin.Context) {
	var category models.Category
	categoryIdStr := c.Param("id")
	categoryName := c.PostForm("category_name")

	categoryID, err := strconv.ParseUint(categoryIdStr, 10, 64)
	if err != nil {
		log.Println("UpdateCategoryAdminHandler: не удалось преобразовать categoryID в uint64", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "не верный формат categoryID"})
		return
	}
	category.ID = categoryID
	category.Name = categoryName

	fileHeader, err := c.FormFile("logoImage")
	if err == nil {
		ext := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("category_%s_%d%s", category.Name, time.Now().Unix(), ext)
		dst := filepath.Join("web", "static", "img", "category", newFilename)
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить изображение"})
			return
		}
		category.Photo = "img/category/" + newFilename
	} else if !errors.Is(err, http.ErrMissingFile) {
		log.Println("UpdateCategoryAdminHandler: ошибка чтения logoImage", err)
	}

	categoryJson, err := json.Marshal(category)
	url := fmt.Sprintf("/api/v1/category/update/%d", category.ID)
	status, _, _, err := api.ProxyTo(c, "http://localhost:8083", "", url, bytes.NewReader(categoryJson))
	if err != nil {
		log.Println("UpdateCategoryAdminHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("UpdateCategoryAdminHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось обновить категорию"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteCategoryAdminHandler(c *gin.Context) {
	categoryID := c.Param("id")

	url := fmt.Sprintf("/api/v1/category/delete/%s", categoryID)
	status, _, _, err := api.ProxyTo(c, "http://localhost:8083", "", url, nil)
	if err != nil {
		log.Println("DeleteCategoryAdminHandler: ошибка вызова Product Service:", err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Product Service недоступен"})
		return
	}

	if status != http.StatusOK {
		log.Println("DeleteCategoryAdminHandler: Product Service вернул статус", status)
		c.AbortWithStatusJSON(status, gin.H{"error": "не удалось удалить категорию"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
