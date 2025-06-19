package storage

import (
	"user_service/internal/db"
	"user_service/internal/models"
)

func SelectUser(userID, role string) (models.User, error) {
	var user models.User
	err := db.DB().Where("id=? and role=?", userID, role).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func SelectUsers(data []uint64, limit int) ([]models.User, error) {
	var users []models.User
	err := db.DB().Limit(limit).Find(&users, data).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	err := db.DB().Find(&users).Error
	return users, err
}

func UpdateRoleUser(userID, updateRole string) error {
	err := db.DB().Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", updateRole).
		Error
	if err != nil {
		return err
	}

	return nil
}
