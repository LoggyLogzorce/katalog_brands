package models

import "time"

type User struct {
	ID        uint64    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (_ *User) TableName() string {
	return "users"
}
