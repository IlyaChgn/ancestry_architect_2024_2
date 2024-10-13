package models

type CreateTreeRequest struct {
	Name string `json:"name"`
}

type TreeResponse struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
}

type AddPermissionRequest struct {
	TreeID uint `json:"treeID"`
	UserID uint `json:"userID"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type Tree struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Layers []Layer `json:"layers"`
}

type Layer struct {
	ID     uint   `json:"id"`
	Number uint   `json:"number"`
	Nodes  []Node `json:"nodes"`
}
