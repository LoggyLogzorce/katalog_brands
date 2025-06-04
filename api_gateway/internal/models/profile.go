package models

import "time"

type Profile struct {
	UserData struct {
		UserID    uint64    `json:"user_id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"user_data"`
	Favorites []struct {
		UserID    uint64    `json:"user_id"`
		ProductID uint64    `json:"product_id"`
		AddedAt   time.Time `json:"added_at"`
	} `json:"favorites"`
	ViewHistory []struct {
		UserId    uint64    `json:"user_id"`
		ProductID uint64    `json:"product_id"`
		ViewedAt  time.Time `json:"viewed_at"`
	} `json:"view_history"`
}

type ProfileResponse struct {
	UserData struct {
		UserID    uint64    `json:"user_id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"user_data"`
	Favorites   []Product `json:"favorites"`
	ViewHistory []Product `json:"view_history"`
}
