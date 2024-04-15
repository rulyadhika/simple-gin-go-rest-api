package validationformatter

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(errs error) string {
	var errString string

	errors := errs.(validator.ValidationErrors)

	for index, err := range errors {
		errString += fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag())

		if index != len(errors)-1 {
			errString += ", "
		}
	}

	return errString
}
