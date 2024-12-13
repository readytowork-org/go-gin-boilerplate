package constants

import (
	"fmt"
)

type UserStatus string // @name UserStatus

const (
	UnVerifiedEmail UserStatus = "unverified-email"
	BaseInfo        UserStatus = "basic-info"
	PlanSelect      UserStatus = "plan-select"
	CardRegister    UserStatus = "card-register"
	CardRegistered  UserStatus = "card-registered"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)

func (r Gender) IsValidVal(val string) error {
	switch val {
	case "male":
		return nil
	case "female":
		return nil
	case "other":
		return nil
	default:
		return fmt.Errorf("invalid user relation status: %q", val)
	}
}
