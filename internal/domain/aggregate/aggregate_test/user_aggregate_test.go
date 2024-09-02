package aggregate_test

import (
	"errors"
	"testing"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUserAggregate(t *testing.T) {
	memData := getUserTestData()
	testCases := []struct {
		name              string
		userModel         model.User
		expectedAggregate *aggregate.UserAggregate
		isError           bool
		errFields         []string
	}{
		{
			name:              "Should correct create aggregate",
			userModel:         memData.correctUser,
			expectedAggregate: &aggregate.UserAggregate{User: memData.correctUser},
			isError:           false,
		},
		{
			name:              "Should required errors",
			userModel:         model.User{},
			expectedAggregate: &aggregate.UserAggregate{User: model.User{}},
			isError:           true,
			errFields:         []string{"Id", "Name", "Value", "Role"},
		},
		{
			name:              "Should uuidv4 error",
			userModel:         memData.incorrectIdUser,
			expectedAggregate: &aggregate.UserAggregate{User: memData.incorrectIdUser},
			isError:           true,
			errFields:         []string{"Id"},
		},
		{
			name:              "Should userRole error",
			userModel:         memData.incorrectUserRole,
			expectedAggregate: &aggregate.UserAggregate{User: memData.incorrectUserRole},
			isError:           true,
			errFields:         []string{"Role"},
		},
		{
			name:              "Should min error",
			userModel:         memData.incorrectMinUser,
			expectedAggregate: &aggregate.UserAggregate{User: memData.incorrectMinUser},
			isError:           true,
			errFields:         []string{"Name"},
		},
		{
			name:              "Should max error",
			userModel:         memData.incorrectMaxUser,
			expectedAggregate: &aggregate.UserAggregate{User: memData.incorrectMaxUser},
			isError:           true,
			errFields:         []string{"Name"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := aggregate.NewUserAggregate(tc.userModel)
			if tc.isError {
				assert.Error(t, err)
				var validationError validator.ValidationErrors
				ok := errors.As(err, &validationError)
				assert.True(t, ok)
				for i, e := range validationError {
					assert.Equal(t, tc.errFields[i], e.Field())
				}
			} else {
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedAggregate, result)
			}
		})
	}
}
