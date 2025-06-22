package storage

import (
	"auth_service/internal/models"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	SelectUser(ctx context.Context, data models.User) (*models.User, error)
	InsertUser(ctx context.Context, data models.User) error
}

type repoUser struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repoUser{db: db}
}

func (r *repoUser) SelectUser(ctx context.Context, data models.User) (*models.User, error) {
	var user *models.User
	err := r.db.WithContext(ctx).Where("email=?", data.Email).First(&user).Error
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}
