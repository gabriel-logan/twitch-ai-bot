package httperrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func BadRequestFromBindError(err error) any {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return "Malformed JSON: " + syntaxErr.Error()
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return "Invalid type for field '" + typeErr.Field + "': expected " + typeErr.Type.String()
	}

	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		return fmt.Sprintf("Invalid value '%s' for field: must be a number", numErr.Num)
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMessages := make([]string, 0, len(validationErrors))

		for _, ve := range validationErrors {
			msg := buildMessage(ve)
			errorMessages = append(errorMessages, msg)
		}

		return errorMessages
	}

	return "Invalid request body: " + err.Error()
}

func buildMessage(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)

	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email", field)

	case "min":
		return fmt.Sprintf("Field '%s' must have at least %s characters", field, fe.Param())

	case "max":
		return fmt.Sprintf("Field '%s' must have at most %s characters", field, fe.Param())

	case "gte":
		return fmt.Sprintf("Field '%s' must be >= %s", field, fe.Param())

	case "lte":
		return fmt.Sprintf("Field '%s' must be <= %s", field, fe.Param())

	default:
		return fmt.Sprintf("Field '%s' is invalid (%s)", field, fe.Tag())
	}
}
