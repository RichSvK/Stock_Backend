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
		field := fe.Field() // already json tag name

		switch fe.Tag() {
		case "required":
			return field + " is required"

		case "oneof":
			return field + " must be one of [" + fe.Param() + "]"

		case "min":
			return field + " must be greater than or equal to " + fe.Param()

		case "max":
			return field + " must be less than or equal to " + fe.Param()

		case "email":
			return field + " must be a valid email"

		default:
			return field + " is invalid"
		}
	}

	return "invalid request"
}
