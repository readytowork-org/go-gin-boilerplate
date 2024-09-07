package user

import (
	"boilerplate-api/api/user/user"
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

// CreateUser to create the CreateUser
func (c Service) CreateUser(user user.CUser) error {
	err := c.repository.Create(user)
	return err
}

// GetAllUsers to get all the CreateUser
func (c Service) GetAllUsers(pagination Pagination) ([]GetUserResponse, int64, error) {
	return c.repository.GetAllUsers(pagination)
}

// GetOneUser one user
func (c Service) GetOneUser(Id int64) (GetUserResponse, error) {
	return c.repository.GetOneUser(Id)
}

// GetOneUserWithEmail Get one user with email
func (c Service) GetOneUserWithEmail(Email string) (user.CUser, error) {
	return c.repository.GetOneUserWithEmail(Email)
}

// GetOneUserWithPhone Get one user with phone
func (c Service) GetOneUserWithPhone(Phone string) (user.CUser, error) {
	return c.repository.GetOneUserWithPhone(Phone)
}
