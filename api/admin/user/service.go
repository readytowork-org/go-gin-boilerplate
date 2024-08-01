package user

import (
	"boilerplate-api/database/models"
	"gorm.io/gorm"
)

type Service struct {
	repository Repository
}

// NewService Creates New user service
func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

// WithTrx repository with transaction
func (c Service) WithTrx(trxHandle *gorm.DB) Service {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

// CreateUser to create the User
func (c Service) CreateUser(user models.User) error {
	err := c.repository.Create(user)
	return err
}

// GetAllUsers to get all the User
func (c Service) GetAllUsers(pagination Pagination) ([]GetUserResponse, int64, error) {
	return c.repository.GetAllUsers(pagination)
}

// GetOneUser one user
func (c Service) GetOneUser(Id int64) (GetUserResponse, error) {
	return c.repository.GetOneUser(Id)
}

// GetOneUserWithEmail Get one user with email
func (c Service) GetOneUserWithEmail(Email string) (models.User, error) {
	return c.repository.GetOneUserWithEmail(Email)
}

// GetOneUserWithPhone Get one user with phone
func (c Service) GetOneUserWithPhone(Phone string) (models.User, error) {
	return c.repository.GetOneUserWithPhone(Phone)
}
