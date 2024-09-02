package appValidator

import "github.com/go-playground/validator/v10"

func isGender(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "male" || value == "female" {
		return true
	}
	return false
}
