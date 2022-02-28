package common

import (
	"github.com/go-playground/validator"
)

var validate = validator.New()

func GetValidate() *validator.Validate {
	return validate
}
