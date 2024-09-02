package auth_usecase_test

import (
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	authUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/auth_usecase"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/google/uuid"
	"net/http"
)

type authUseCaseRegistrationDataMock struct {
	CorrectRegInput1          appDto.RegistrationUseCaseDto
	CorrectRegInput1Result1   *authUseCase.AuthResult
	CorrectRegInput1Result2   *appErrors.AppError
	IncorrectRegInput2        appDto.RegistrationUseCaseDto
	IncorrectRegInput2Result1 *authUseCase.AuthResult
	IncorrectRegInput2Result2 *appErrors.AppError
}

type authUseCaseLoginDataMock struct {
	CorrectRegInput1          appDto.LoginUseCaseDto
	CorrectRegInput1Result1   *authUseCase.AuthResult
	CorrectRegInput1Result2   *appErrors.AppError
	IncorrectRegInput1        appDto.LoginUseCaseDto
	IncorrectRegInput1Result1 *authUseCase.AuthResult
	IncorrectRegInput1Result2 *appErrors.AppError
	IncorrectRegInput2        appDto.LoginUseCaseDto
	IncorrectRegInput2Result1 *authUseCase.AuthResult
	IncorrectRegInput2Result2 *appErrors.AppError
}

type authUseCaseDataMock struct {
	Registration *authUseCaseRegistrationDataMock
	Login        *authUseCaseLoginDataMock
}

func newAuthUseCaseDataMock() *authUseCaseDataMock {
	return &authUseCaseDataMock{
		Registration: &authUseCaseRegistrationDataMock{
			CorrectRegInput1: appDto.RegistrationUseCaseDto{
				Name:     "NewEer0",
				Password: "Supperpupper123",
			},
			CorrectRegInput1Result1: &authUseCase.AuthResult{
				User: &appDto.ResponseUserDto{
					Id:   uuid.New().String(),
					Name: "NewEer0",
					Role: constants.UserRole,
				},
				Tokens: tokenService.JwtTokens{AccessToken: "dsads", RefreshToken: "dsadsad"},
			},
			CorrectRegInput1Result2: nil,
			IncorrectRegInput2: appDto.RegistrationUseCaseDto{
				Name:     "NewEer02",
				Password: "incorrect_passord",
			},
			IncorrectRegInput2Result1: nil,
			IncorrectRegInput2Result2: &appErrors.AppError{
				Code:    http.StatusUnprocessableEntity,
				Message: appErrors.DefaultUnprocessableEntity,
			},
		},
		Login: &authUseCaseLoginDataMock{
			CorrectRegInput1: appDto.LoginUseCaseDto{
				Name:     "NewEer0",
				Password: "Supperpupper123",
			},
			CorrectRegInput1Result1: &authUseCase.AuthResult{
				User: &appDto.ResponseUserDto{
					Id:   uuid.New().String(),
					Name: "NewEer0",
					Role: constants.UserRole,
				},
				Tokens: tokenService.JwtTokens{AccessToken: "dsads", RefreshToken: "dsadsad"},
			},
			CorrectRegInput1Result2: nil,
			IncorrectRegInput2: appDto.LoginUseCaseDto{
				Name:     "NewEer0",
				Password: "incorrect_passord",
			},
			IncorrectRegInput2Result1: nil,
			IncorrectRegInput2Result2: &appErrors.AppError{
				Code:    http.StatusForbidden,
				Message: constants.NickOrPasswordIncorrect,
			},
			IncorrectRegInput1: appDto.LoginUseCaseDto{
				Name:     "NotUser",
				Password: "Supperpupper123",
			},
			IncorrectRegInput1Result1: nil,
			IncorrectRegInput1Result2: &appErrors.AppError{
				Code:    http.StatusForbidden,
				Message: constants.NickOrPasswordIncorrect,
			},
		},
	}
}
