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

// GetOneUser one user
func (c Service) GetOneUser(Id string) (models.User, error) {
	return c.repository.GetOneUser(Id)
}
