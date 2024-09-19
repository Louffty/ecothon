package dto

// CreateUser is a struct that represents a DTO to create a new User.
type CreateUser struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

// AuthUser is a struct that represents a DTO to authenticate a User.
type AuthUser struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required,password"`
}

type Author struct {
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	Rate       string `json:"rate"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}

type ReturnUser struct {
	UUID        string            `json:"uuid"`
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	Role        string            `json:"role"`
	CoinsAmount int               `json:"coins_amount"`
	Rate        string            `json:"rate"`
	Events      ReturnUsersEvents `json:"events"`
	IsVerified  bool              `json:"is_verified"`
}

type VerifiedUser struct {
	Organisation string `json:"organisation"`
}
