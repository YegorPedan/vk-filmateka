package token_servicetest_test

import (
	"fmt"
	"testing"

	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestTokenServiceGenerate(t *testing.T) {
	cfg := config.MustLoad()
	tokenRepo := mockRepository.NewTokenRepository()
	tokenServ := tokenService.New(tokenRepo)

	jwtData := tokenService.JwtUserData{Id: "my-uuidv4", Role: constants.UserRole}
	tokens, err := tokenServ.Generate(jwtData)
	assert.Nil(t, err)
	refreshToken, err := jwt.Parse(tokens.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.ApiKey), nil
	})
	assert.True(t, refreshToken.Valid)
}
