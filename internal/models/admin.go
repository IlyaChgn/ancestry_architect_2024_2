package models

import "time"

type AdminResponse struct {
	ID     uint   `json:"id"`
	IsAuth bool   `json:"isAuth"`
	Email  string `json:"email"`
}

type EditPasswordRequest struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

type UserForAdminResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type UsersList []UserForAdminResponse

type EditTreeRequest struct {
	TreeID uint   `json:"treeID"`
	Name   string `json:"name"`
}

type NodeForAdmin struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Birthdate   time.Time `json:"birthdate"`
	Deathdate   time.Time `json:"deathdate"`
	Gender      string    `json:"gender"`
	PreviewPath string    `json:"previewPath"`
	LayerID     uint      `json:"layerID"`
	LayerNum    int       `json:"layerNum"`
	TreeID      uint      `json:"treeID"`
	UserID      uint      `json:"userID"`
	IsDeleted   bool      `json:"isDeleted"`
}
