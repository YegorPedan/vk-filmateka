package tokenService

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	"github.com/golang-jwt/jwt"
)

type (
	JwtTokens struct {
		AccessToken  string
		RefreshToken string
	}

	JwtUserData struct {
		Id   string `json:"id"`
		Role string `json:"role"`
	}

	CustomClaims struct {
		JwtUserData `json:"jwtUserData"`
		jwt.StandardClaims
	}

	Service interface {
		HasByValue(ctx context.Context, refreshToken string) (bool, error)
		Generate(data JwtUserData) (*JwtTokens, error)
		ValidateRefreshToken(refreshToken string) (*JwtUserData, error)
		Save(ctx context.Context, data appDto.SaveTokenServiceDto) (*model.Token, error)
		DeleteByValue(ctx context.Context, value string) error
	}

	tokenService struct {
		repository.TokenRepository
	}
)

func New(tokenRepo repository.TokenRepository) Service {
	return &tokenService{
		TokenRepository: tokenRepo,
	}
}
