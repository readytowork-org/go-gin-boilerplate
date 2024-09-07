package user

import (
	"boilerplate-api/internal/config"
	"gorm.io/gorm"
)

// Repository database structure
type Repository struct {
	db     *config.Database
	logger config.Logger
}

// NewRepository creates a new CUser repository
func NewRepository(db *config.Database, logger config.Logger) Repository {
	return Repository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c Repository) WithTrx(trxHandle *gorm.DB) Repository {
	if trxHandle == nil {
		c.logger.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

func (c Repository) GetOneUser(Id string) (userModel CUser, err error) {
	return userModel, c.db.DB.
		Model(&userModel).
		Where("id = ?", Id).
		First(&userModel).
		Error
}
