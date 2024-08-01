package types

import (
	"time"

	"gorm.io/gorm"
)

type TimeStamps struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` //add soft delete in gorm
}

// BaseModal contains common columns for all tables.
type BaseModal struct {
	ID int64 `json:"id"`
	TimeStamps
}
