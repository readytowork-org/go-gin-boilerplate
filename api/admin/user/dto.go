package user

import (
	"boilerplate-api/api/user/user"
)

// CreateUserRequestData Request body data to create user
type CreateUserRequestData struct {
	user.CUser
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// GetUserResponse Dtos for CUser model
type GetUserResponse struct {
	user.CUser
	Password string
}
