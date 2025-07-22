package common

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatorEmail(f1 validator.FieldLevel) bool {
	email := f1.Field().String()
	ok, _ := regexp.MatchString(`^\w{5,}@[a-z0-9]{2,3}\.[a-z]+$|\,$`, email)
	return ok
}
