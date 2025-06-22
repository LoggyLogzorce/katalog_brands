package storage

import (
	"auth_service/internal/models"
	"context"
	"fmt"
)

func (r *repoUser) InsertUser(ctx context.Context, data models.User) error {
	if data.Role == "" {
		data.Role = "user"
	}
	err := r.db.WithContext(ctx).Save(&data).Error
	fmt.Println(err)
	return err
}
