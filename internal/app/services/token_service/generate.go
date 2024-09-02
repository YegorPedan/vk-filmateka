package tokenService

import (
	"time"

	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	"github.com/golang-jwt/jwt"
)

func (t *tokenService) Generate(data JwtUserData) (*JwtTokens, error) {
	cfg := config.NewConfig()
	accessDuration, err := time.ParseDuration(cfg.AccessTokenTime)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: TokenService, method: Generate. ", "Parse access token time duration")
	}
	refreshDuration, err := time.ParseDuration(cfg.RefreshTokenTime)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: TokenService, method: Generate. ", "Parse refresh token time duration")
	}
	accessClaims := CustomClaims{
		JwtUserData:    data,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(accessDuration).Unix()},
	}
	refreshClaims := CustomClaims{
		JwtUserData:    data,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(refreshDuration).Unix()},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.ApiKey))
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: TokenService, method: Generate. ", "Sign access token error: ", err.Error())
	}
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.ApiKey))
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: TokenService, method: Generate. ", "Sign refresh token error: ", err.Error())
	}
	return &JwtTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
