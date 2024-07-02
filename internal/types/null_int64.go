package types

import "database/sql/driver"

// NullInt64 represents an int64 that may be null.
// NullInt64 implements the [Scanner] interface so
// it can be used as a scan destination, similar to [NullString].
type NullInt64 struct {
	Int64 int64 `json:"Int64"`
	Valid bool  `json:"Valid"`
} // @name NullInt64

// Scan implements the [Scanner] interface.
func (n *NullInt64) Scan(value any) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	n.Int64 = value.(int64)

	return nil
}

// Value implements the [driver.Valuer] interface.
func (n *NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}
