package types

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// BINARY16 -> new datatype
type BINARY16 uuid.UUID

// StringToBinary16 -> parse string to id
func StringToBinary16(s string) (BINARY16, error) {
	id, err := uuid.Parse(s)
	return BINARY16(id), err
}

// String -> String Representation of Binary16
func (binary *BINARY16) String() string {
	return uuid.UUID(*binary).String()
}

// GormDataType -> sets type to binary(16)
func (binary *BINARY16) GormDataType() string {
	return "binary(16)"
}

// Scan --> From DB
func (binary *BINARY16) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal  value:", value))
	}
	parseByte, err := uuid.FromBytes(bytes)
	*binary = BINARY16(parseByte)
	return err
}

// Value -> TO DB
func (binary *BINARY16) Value() (driver.Value, error) {
	return uuid.UUID(*binary).MarshalBinary()
}

// MarshalJSON -> convert to json string
func (binary *BINARY16) MarshalJSON() ([]byte, error) {
	s := uuid.UUID(*binary)
	str := "\"" + s.String() + "\""
	return []byte(str), nil
}

// UnmarshalJSON -> convert from json string
func (binary *BINARY16) UnmarshalJSON(by []byte) error {
	s, err := uuid.ParseBytes(by)
	*binary = BINARY16(s)
	return err
}
