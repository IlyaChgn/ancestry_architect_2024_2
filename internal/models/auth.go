package models

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupRequest struct {
	UserLoginRequest
	PasswordRepeat string `json:"password_repeat"`
}

type UserResponse struct {
	IsAuth  bool   `json:"is_auth"`
	User    User   `json:"user"`
	Name    string `json:"name"`
	Surname string `json:"surnname"`
}
