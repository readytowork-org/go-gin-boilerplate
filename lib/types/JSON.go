package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JsonType []string

func (jsonType *JsonType) GormDataType() string {
	return "json"
}
func (jsonType *JsonType) Value() (driver.Value, error) {
	return json.Marshal(jsonType)
}
func (jsonType *JsonType) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	var result JsonType
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return errors.New("failed to unmarshal json value")
	}
	*jsonType = result
	return nil
}
