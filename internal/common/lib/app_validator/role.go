package appValidator

import "github.com/go-playground/validator/v10"

func userRole(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	correctRoles := []string{"ADMIN", "USER"}
	for _, item := range correctRoles {
		if item == value {
			return true
		}
	}
	return false
}
