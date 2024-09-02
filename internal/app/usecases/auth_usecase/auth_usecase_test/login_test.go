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
	"testing"
)

func TestAuthLogin(t *testing.T) {
	mockData := newAuthUseCaseDataMock()

	testCases := []struct {
		name           string
		inputData      appDto.LoginUseCaseDto
		expectedResult *authUseCase.AuthResult
		expectedError  *appErrors.AppError
		isError        bool
	}{
		{
			name:           "Shoul correct login user",
			inputData:      mockData.Login.CorrectRegInput1,
			expectedResult: mockData.Login.CorrectRegInput1Result1,
			expectedError:  mockData.Login.CorrectRegInput1Result2,
			isError:        false,
		},
		{
			name:           "Shoul incorrect login password",
			inputData:      mockData.Login.IncorrectRegInput2,
			expectedResult: mockData.Login.IncorrectRegInput2Result1,
			expectedError:  mockData.Login.IncorrectRegInput2Result2,
			isError:        true,
		},
		{
			name:           "Shoul incorrect login name",
			inputData:      mockData.Login.IncorrectRegInput1,
			expectedResult: mockData.Login.IncorrectRegInput1Result1,
			expectedError:  mockData.Login.IncorrectRegInput1Result2,
			isError:        true,
		},
	}

	userRepo := mockRepository.NewUserRepository()
	tokenRepo := mockRepository.NewTokenRepository()

	useCase := authUseCase.New(userService.New(userRepo), tokenService.New(tokenRepo), userRepo)
	cfg := config.MustLoad()
	_, _ = useCase.Registration(context.Background(), mockData.Registration.CorrectRegInput1)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := useCase.Login(context.Background(), testCase.inputData)
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
