package auth_usecase_test

import (
	"context"
	"errors"
	"fmt"
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	userService "github.com/OddEer0/vk-filmoteka/internal/app/services/user_service"
	authUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/auth_usecase"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func isEqualUser(a *appDto.ResponseUserDto, b *appDto.ResponseUserDto) bool {
	aType := reflect.TypeOf(*a)
	aValue := reflect.ValueOf(*a)
	bValue := reflect.ValueOf(*b)
	for i := 0; i < aType.NumField(); i++ {
		aField := aType.Field(i)
		aVal := aValue.Field(i).Interface()
		bVal := bValue.Field(i).Interface()
		if aField.Name == "Id" {
			continue
		}
		if aVal != bVal {
			return false
		}
	}
	return true
}

func TestAuthRegistration(t *testing.T) {
	mockData := newAuthUseCaseDataMock().Registration

	testCases := []struct {
		name           string
		inputData      appDto.RegistrationUseCaseDto
		expectedResult *authUseCase.AuthResult
		expectedError  *appErrors.AppError
		isError        bool
	}{
		{
			name:           "Should register user",
			inputData:      mockData.CorrectRegInput1,
			expectedResult: mockData.CorrectRegInput1Result1,
			expectedError:  mockData.CorrectRegInput1Result2,
			isError:        false,
		},
		{
			name:           "Should incorrect password unproccessable error",
			inputData:      mockData.IncorrectRegInput2,
			expectedResult: mockData.IncorrectRegInput2Result1,
			expectedError:  mockData.IncorrectRegInput2Result2,
			isError:        true,
		},
	}

	userRepo := mockRepository.NewUserRepository()
	tokenRepo := mockRepository.NewTokenRepository()

	useCase := authUseCase.New(userService.New(userRepo), tokenService.New(tokenRepo), userRepo)
	cfg := config.MustLoad()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := useCase.Registration(context.Background(), testCase.inputData)
			if !testCase.isError {
				assert.Nil(t, err)
				refreshToken, err := jwt.Parse(result.Tokens.RefreshToken, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(cfg.ApiKey), nil
				})
				assert.Equal(t, nil, err)
				assert.True(t, refreshToken.Valid)
				assert.NotEmpty(t, result.Tokens.RefreshToken)
				assert.NotEmpty(t, result.Tokens.AccessToken)
				assert.True(t, isEqualUser(testCase.expectedResult.User, result.User))
			} else {
				assert.Nil(t, testCase.expectedResult)
				var appErr *appErrors.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, testCase.expectedError.Code, appErr.Code)
					assert.Equal(t, testCase.expectedError.Message, appErr.Message)
				} else {
					t.Error("incorrect type error!!!")
				}
			}
		})
	}

	db := inMemDb.New()
	db.CleanUp()
}
