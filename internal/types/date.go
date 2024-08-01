package types

import (
	"database/sql/driver"
	"time"
)

type Date string //	@name Date

// GormDataType -> sets type to binary(16)
func (date *Date) GormDataType() string {
	return "DATE"
}

// Scan implements the [Scanner] interface.
func (date *Date) Scan(value any) error {
	_date := value.(time.Time)
	*date = Date(_date.Format("2006-01-02"))
	return nil
}

// Value implements the [driver.Valuer] interface.
func (date *Date) Value() (driver.Value, error) {
	return date, nil
}
