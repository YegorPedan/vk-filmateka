package vo_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/OddEer0/vk-filmoteka/internal/domain/valuesobject"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordVO(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		result   valuesobject.Password
		err      error
	}{
		{
			name:     "Should correct create password",
			password: "imsupervalid12",
			err:      nil,
			result:   valuesobject.Password{Value: "imsupervalid12"},
		},
		{
			name:     "Should incorrect create password",
			password: "imincorrect",
			err:      errors.New(valuesobject.PasswordInvalid),
			result:   valuesobject.Password{},
		},
		{
			name:     "Should incorrect create password 2",
			password: "123456789",
			err:      errors.New(valuesobject.PasswordInvalid),
			result:   valuesobject.Password{},
		},
		{
			name:     "Should incorrect min 8 password",
			password: "cor2",
			err:      errors.New(valuesobject.PasswordMinLength),
			result:   valuesobject.Password{},
		},
		{
			name:     "Should incorrect create password 2",
			password: strings.Repeat("pepes", 7) + "1",
			err:      errors.New(valuesobject.PasswordMaxLength),
			result:   valuesobject.Password{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pass, err := valuesobject.NewPassword(tc.password)
			assert.Equal(t, err, tc.err)
			if tc.err == nil {
				err2 := bcrypt.CompareHashAndPassword([]byte(pass.Value), []byte(tc.password))
				assert.Equal(t, nil, err2)
			} else {
				assert.Equal(t, pass, valuesobject.Password{})
			}
		})
	}
}
