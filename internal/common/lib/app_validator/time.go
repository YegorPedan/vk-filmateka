package appValidator

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func dateIsLessNow(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return !date.IsZero() && date.Before(time.Now())
}
