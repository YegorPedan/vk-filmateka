package valuesobject

import (
	"encoding/json"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordMinLength = "password min length 8"
	PasswordMaxLength = "password max length 35"
	PasswordInvalid   = "invalid password"
)

type Password struct {
	Value string `validate:"required"`
}

func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

// IsValidPassword TODO - Пофиксить проблему с регулркой
func (p Password) IsValidPassword(password string) bool {
	re := regexp.MustCompile(`[[:digit:]]`)
	lowercaseCheck := regexp.MustCompile(`[[:lower:]]`).MatchString(password)
	return re.MatchString(password) && lowercaseCheck
}

func (p Password) Validate(password string) error {
	length := len(password)
	switch {
	case length < 8:
		return errors.New(PasswordMinLength)
	case length > 35:
		return errors.New(PasswordMaxLength)
	case !p.IsValidPassword(password):
		return errors.New(PasswordInvalid)
	}
	return nil
}

func NewPassword(password string) (Password, error) {
	result := Password{Value: password}
	err := result.Validate(password)
	if err != nil {
		return Password{}, err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	result.Value = string(hash)
	return result, nil
}
