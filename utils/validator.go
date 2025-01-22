package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// AmountValidator checks if the amount field has to be string and up to 2 decimal places are sent.
func AmountValidator(fl validator.FieldLevel) bool {
	amount := fl.Field().String()
	return regexp.MustCompile(`^\d+(\.\d{1,2})?$`).MatchString(amount)
}
