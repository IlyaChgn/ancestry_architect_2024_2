package models

import (
	"github.com/jackc/pgtype"
)

type CreateNodeRequest struct {
	IsFirstNode bool             `json:"isFirstNode"`
	TreeID      uint             `json:"treeID"`
	Name        string           `json:"name"`
	Addition    AdditionDataList `json:"addition"`
	Relatives   GetRelativesList `json:"relatives"`
	IsSpouse    bool             `json:"isSpouse"`
	Gender      string           `json:"gender"`
}

type GetRelativesList struct {
	Children []uint `json:"children"`
	Parents  []uint `json:"parents"`
	Spouses  []uint `json:"spouses"`
	Siblings []uint `json:"siblings"`
}

type AdditionDataList struct {
	Birthdate   string `json:"birthdate"`
	Deathdate   string `json:"deathdate"`
	Description string `json:"description"`
}

type Node struct {
	ID          uint              `json:"id"`
	LayerID     uint              `json:"layerID"`
	Name        string            `json:"name"`
	Birthdate   *pgtype.Date      `json:"birthdate"`
	Deathdate   *pgtype.Date      `json:"deathdate"`
	PreviewPath string            `json:"previewPath"`
	Relatives   SendRelativesList `json:"relatives"`
	IsSpouse    bool              `json:"isSpouse"`
	Gender      string            `json:"gender"`
}

type SendRelativesList struct {
	Parents  []uint `json:"parents"`
	Spouses  []uint `json:"spouses"`
	Children []uint `json:"children"`
}

type DescriptionResponse struct {
	NodeID      uint   `json:"nodeID"`
	Description string `json:"description"`
}

type EditNodeRequest struct {
	Name        string `json:"name"`
	Birthdate   string `json:"birthdate"`
	Deathdate   string `json:"deathdate"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
}

type EditNodeResponse struct {
	Name      string       `json:"name"`
	Birthdate *pgtype.Date `json:"birthdate"`
	Deathdate *pgtype.Date `json:"deathdate"`
	Gender    string       `json:"gender"`
}

type UpdatePreviewResponse struct {
	ID          uint   `json:"id"`
	PreviewPath string `json:"previewPath"`
}

type Relative struct {
	ID          uint
	LayerNumber int
	LayerID     uint
}

type ReturningAdditionalData struct {
	Birthdate *pgtype.Date
	Deathdate *pgtype.Date
}
