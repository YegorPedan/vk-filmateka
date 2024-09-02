package authUseCase

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	userService "github.com/OddEer0/vk-filmoteka/internal/app/services/user_service"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
)

type (
	AuthResult struct {
		User   *appDto.ResponseUserDto
		Tokens tokenService.JwtTokens
	}

	AuthUseCase interface {
		Registration(ctx context.Context, data appDto.RegistrationUseCaseDto) (*AuthResult, error)
		Login(ctx context.Context, data appDto.LoginUseCaseDto) (*AuthResult, error)
		Logout(ctx context.Context, refreshToken string) error
		Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
	}

	authUseCase struct {
		repository.UserRepository
		UserService  userService.Service
		TokenService tokenService.Service
	}
)

func New(userService userService.Service, tokenService tokenService.Service, userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{UserService: userService, TokenService: tokenService, UserRepository: userRepo}
}
