package user

import (
	"boilerplate-api/database/models"
)

// CreateUserRequestData Request body data to create user
type CreateUserRequestData struct {
	models.User
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// GetUserResponse Dtos for User model
type GetUserResponse struct {
	models.User
	Password string
}
