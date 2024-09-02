package tokenService

import (
	"errors"

	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	"github.com/golang-jwt/jwt"
)

func (t *tokenService) ValidateRefreshToken(refreshToken string) (*JwtUserData, error) {
	cfg := config.NewConfig()
	token, err := jwt.ParseWithClaims(refreshToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.ApiKey), nil
	})
	if err != nil {
		var jwtErr *jwt.ValidationError
		if errors.As(err, &jwtErr) {
			return nil, jwtErrHandle(jwtErr)
		}
		return nil, appErrors.InternalServerError("", "target: TokenService, method: ValidateRefreshToken. ", "jwt parse error: ", err.Error())
	}
	if !token.Valid {
		return nil, appErrors.Unauthorized("")
	}
	claim := token.Claims.(*CustomClaims)
	return &claim.JwtUserData, nil
}

func jwtErrHandle(jwtErr *jwt.ValidationError) error {
	if jwtErr.Errors&jwt.ValidationErrorMalformed != 0 {
		return appErrors.Unauthorized("", "target: TokenService. ", "Uncorrected jwt token")
	} else if jwtErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
		return appErrors.Unauthorized("", "target: TokenService. ", "Токен истек или еще не начал действовать")
	}
	return appErrors.Unauthorized("", "target: TokenService. ", "Ошибка проверки подписи токена")
}
