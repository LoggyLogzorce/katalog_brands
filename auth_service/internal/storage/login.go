package storage

import (
	"auth_service/internal/db"
	"auth_service/internal/models"
)

func SelectUser(data models.User) (*models.User, error) {
	var user *models.User
	err := db.DB().Where("email=?", data.Email).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}
