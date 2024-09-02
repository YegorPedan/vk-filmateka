package user_service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	userService "github.com/OddEer0/vk-filmoteka/internal/app/services/user_service"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceCreate(t *testing.T) {
	testCases := []struct {
		name    string
		newUser appDto.RegistrationUseCaseDto
		isError bool
		errCode int
	}{
		{
			name: "Should create user aggregate",
			newUser: appDto.RegistrationUseCaseDto{
				Name:     "new_user",
				Password: "CorrectPassword123",
			},
			isError: false,
		},
		{
			name: "Should conflict error by Name user",
			newUser: appDto.RegistrationUseCaseDto{
				Name:     "new_user",
				Password: "CorrectPassword123",
			},
			isError: true,
			errCode: http.StatusConflict,
		},
		{
			name: "Should error incorrect password",
			newUser: appDto.RegistrationUseCaseDto{
				Name:     "new_user23",
				Password: "incorrect",
			},
			isError: true,
			errCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Should error incorrect Name",
			newUser: appDto.RegistrationUseCaseDto{
				Name:     "new_user23dsadsadsadsadsagsadsadsadsadsadsadsadsads",
				Password: "incorrect",
			},
			isError: true,
			errCode: http.StatusUnprocessableEntity,
		},
	}

	userRepo := mockRepository.NewUserRepository()
	userServ := userService.New(userRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := userServ.Create(context.Background(), tc.newUser)
			if !tc.isError {
				if assert.Equal(t, nil, err) {
					assert.NotEmpty(t, result.User.Id)
					assert.Equal(t, constants.UserRole, result.User.Role)
					assert.Equal(t, tc.newUser.Name, result.User.Name)
				}
				userRepo.Create(context.Background(), result)
			} else {
				assert.Equal(t, (*aggregate.UserAggregate)(nil), result)
				var appErr *appErrors.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tc.errCode, appErr.Code)
				}
			}
		})
	}
}
