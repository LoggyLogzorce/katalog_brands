package storage

import (
	"auth_service/internal/db"
	"auth_service/internal/models"
	"fmt"
)

func InsertUser(data models.User) error {
	if data.Role == "" {
		data.Role = "user"
	}
	err := db.DB().Save(&data).Error
	fmt.Println(err)
	return err
}
