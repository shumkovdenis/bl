package helpers

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidateError struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func Validate(s interface{}) []*ValidateError {
	var errors []*ValidateError
	if err := validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidateError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
