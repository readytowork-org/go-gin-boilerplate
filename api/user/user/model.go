package user

import (
	"boilerplate-api/database/dao"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CUser struct {
	dao.User
}

// BeforeCreate Runs before inserting a row into table
func (u *CUser) BeforeCreate(db *gorm.DB) error {
	var Zap *zap.SugaredLogger
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(password)
	if err != nil {
		Zap.Error("Error decrypting plain password to hash", err.Error())
		return err
	}
	return nil
}
