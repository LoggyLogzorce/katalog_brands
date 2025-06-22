package storage

import (
	"context"
	"gorm.io/gorm"
	"user_service/internal/models"
)

type UserRepository interface {
	SelectUser(ctx context.Context, userID, role string) (models.User, error)
	SelectUsers(ctx context.Context, data []uint64, limit int) ([]models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	UpdateRoleUser(ctx context.Context, userID, updateRole string) error
	SelectHistory(ctx context.Context, userID string, limit int) ([]models.History, error)
	DeleteViewHistory(ctx context.Context, userID, productID string) error
	CreateView(ctx context.Context, userID, productID uint64) error
	SelectView(ctx context.Context, userID, productID uint64) (bool, error)
	CountView(ctx context.Context, data []uint64) (int, error)
	GetProductViewsStats(ctx context.Context, productsID []uint64) ([]models.ProductViewStat, error)
	SelectFavorite(ctx context.Context, userID string, limit int) ([]models.Favorite, error)
	CreateFavorite(ctx context.Context, userID, productID uint64) error
	DeleteFavorite(ctx context.Context, userID, productID string) error
}

type repoUser struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repoUser{db: db}
}

func (r *repoUser) SelectUser(ctx context.Context, userID, role string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id=? and role=?", userID, role).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *repoUser) SelectUsers(ctx context.Context, data []uint64, limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Limit(limit).Find(&users, data).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repoUser) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *repoUser) UpdateRoleUser(ctx context.Context, userID, updateRole string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", updateRole).
		Error
	if err != nil {
		return err
	}

	return nil
}
