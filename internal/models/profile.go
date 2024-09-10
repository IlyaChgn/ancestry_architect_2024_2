package models

import (
	"mime/multipart"
	"time"

	"github.com/jackc/pgtype"
)

type Profile struct {
	ID         uint        `json:"id"`
	UserID     uint        `json:"userID"`
	Name       string      `json:"name"`
	Surname    string      `json:"surname"`
	Birthdate  pgtype.Date `json:"birthdate"`
	Gender     string      `json:"gender"`
	AvatarPath string      `json:"avatarPath"`
}

type ProfileNullData struct {
	Name       *string
	Surname    *string
	AvatarPath *string
	Gender     *string
}

type UpdateProfileRequest struct {
	Email     string
	Name      string
	Surname   string
	Birthdate time.Time
	Gender    string
	Avatar    *multipart.FileHeader
}

type UpdateProfileResponse struct {
	User    User    `json:"user"`
	Profile Profile `json:"profile"`
}
