package helper

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Validation Error formats validator errors
func ValidationError(err error) string {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		fe := ve[0]
		field := fe.Field()

		switch fe.Tag() {
		case "min":
			return field + " must be greater than or equal to " + fe.Param()

		case "required":
			return field + " is required"

		case "oneof":
			return field + " must be one of [" + fe.Param() + "]"

		case "max":
			return field + " must be less than or equal to " + fe.Param()
		
		case "len":
			return field + " must be exactly " + fe.Param() + " characters long"
			
		case "alpha":
			return field + " must contain only alphabetic characters"

		default:
			return field + " is invalid"
		}
	}

	return "invalid request"
}
