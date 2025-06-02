package models

import "time"

type User struct {
	ID           uint64    `json:"userID"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (_ *User) TableName() string {
	return "users"
}
