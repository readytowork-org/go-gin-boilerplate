package request_validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"boilerplate-api/lib/api_errors"
	"boilerplate-api/lib/constants"

	"github.com/go-playground/validator/v10"
)

// FIXME :: refactor validator | use one validator for all

// Validator structure
type Validator struct {
	*validator.Validate
}

// NewValidator Register Custom Validators
func NewValidator() Validator {
	v := validator.New()
	v.RegisterValidation(
		"phone", func(fl validator.FieldLevel) bool {
			if fl.Field().String() != "" {
				match, _ := regexp.MatchString("^[- +()]*[0-9][- +()0-9]*$", fl.Field().String())
				return match
			}
			return true
		},
	)
	v.RegisterValidation(
		"gender", func(fl validator.FieldLevel) bool {
			if fl.Field().String() != "" {
				var valType constants.Gender
				if err := valType.IsValidVal(fl.Field().String()); err != nil {
					return false
				}
			}
			return true
		},
	)
	v.RegisterValidation(
		"email", func(fl validator.FieldLevel) bool {
			if fl.Field().String() != "" {
				match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, fl.Field().String())
				return match
			}
			return true
		},
	)
	v.RegisterValidation(
		"date", func(fl validator.FieldLevel) bool {
			_, parseErr := time.Parse("2006-01-02", fl.Field().String())
			return parseErr == nil
		},
	)
	v.RegisterValidation(
		"required_if", func(fl validator.FieldLevel) bool {
			//  expected tag format
			// "required_if=OtherField Value1 Value2"
			params := strings.Split(fl.Param(), " ")
			paramsSize := len(params)
			if paramsSize < 2 {
				return false
			}
			paramField := params[0]
			otherFieldValue := fl.Parent().FieldByName(paramField)

			expectedValuesArr := params[1:paramsSize]

			var expectedValueInStr string
			switch otherFieldValue.Kind() {
			case reflect.String:
				expectedValueInStr = otherFieldValue.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				expectedValueInStr = strconv.FormatInt(otherFieldValue.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				expectedValueInStr = strconv.FormatUint(otherFieldValue.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				expectedValueInStr = strconv.FormatFloat(otherFieldValue.Float(), 'f', -1, 64)
			default:
				return true
			}

			for _, expectedValue := range expectedValuesArr {
				if expectedValueInStr == expectedValue && fl.Field().String() == "" {
					return false
				}
			}

			return true
		},
	)
	return Validator{
		Validate: v,
	}
}

func (cv Validator) generateValidationMessage(field string, rule string, param string) (message string) {
	switch rule {
	case "required":
		return fmt.Sprintf("Field '%s' is '%s'.", field, rule)
	case "phone":
		return fmt.Sprintf("Field '%s' is not valid.", field)
	case "gender":
		return fmt.Sprintf("Field '%s' is not valid.", field)
	case "email":
		return fmt.Sprintf("Field '%s' is not valid.", field)
	case "date":
		return fmt.Sprintf("Invalid date format for '%s' use YYYY-MM-DD", field)
	case "required_if":
		params := strings.Split(param, " ")
		requiredIfParent := params[0]
		values := ""
		paramsLen := len(params)
		if paramsLen > 1 {
			values = strings.Join(params[1:paramsLen], ",")
		}
		return fmt.Sprintf("Field %s is required when %s is %s", field, requiredIfParent, values)
	default:
		return fmt.Sprintf("Field '%s' is not valid.", field)
	}
}

func (cv Validator) GenerateValidationResponse(err error) []api_errors.ValidationError {
	var validations []api_errors.ValidationError
	for _, value := range err.(validator.ValidationErrors) {
		field, rule := value.Field(), value.Tag()
		validation := api_errors.ValidationError{
			Field: field, Message: cv.generateValidationMessage(field, rule, value.Param()),
		}
		validations = append(validations, validation)
	}
	return validations
}
