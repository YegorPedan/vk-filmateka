package appValidator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var re *regexp.Regexp = regexp.MustCompile(`[[:digit:]]`)
var lowercase *regexp.Regexp = regexp.MustCompile(`[[:lower:]]`)

func isValidPassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return re.MatchString(value) && lowercase.MatchString(value)
}
