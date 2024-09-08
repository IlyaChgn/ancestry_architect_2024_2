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
	PasswordRepeat string `json:"passwordRepeat"`
}

type UserResponse struct {
	IsAuth  bool   `json:"isAuth"`
	User    User   `json:"user"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
