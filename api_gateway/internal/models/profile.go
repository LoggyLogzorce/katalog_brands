package models

type Profile struct {
	UserData    UserData      `json:"user_data"`
	Favorites   []Favorite    `json:"favorites"`
	ViewHistory []ViewHistory `json:"view_history"`
}

type ProfileResponse struct {
	UserData    UserData  `json:"user_data"`
	Favorites   []Product `json:"favorites"`
	ViewHistory []Product `json:"view_history"`
}
