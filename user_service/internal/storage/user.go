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

func UpdateRoleUser(userID, role, updateRole string) error {
	err := db.DB().Model(&models.User{}).
		Where("id = ? and role=?", userID, role).
		Update("role", updateRole).
		Error
	if err != nil {
		return err
	}

	return nil
}
