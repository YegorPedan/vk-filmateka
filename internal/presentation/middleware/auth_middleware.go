package middleware

import (
	"context"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	"github.com/golang-jwt/jwt"
	"net/http"
	"slices"
)

func AuthRoleMiddleware(roles ...string) func(appErrors.AppHandlerFunc) appErrors.AppHandlerFunc {
	return func(next appErrors.AppHandlerFunc) appErrors.AppHandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) error {
			cfg := config.NewConfig()
			if cfg.Env != "test" {
				accessToken, err := req.Cookie("accessToken")
				if err != nil {
					return appErrors.Unauthorized("")
				}

				token, err := jwt.ParseWithClaims(accessToken.Value, &tokenService.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.ApiKey), nil
				})

				if err != nil {
					return appErrors.Unauthorized("")
				}
				userData := token.Claims.(*tokenService.CustomClaims).JwtUserData

				if !slices.Contains(roles, userData.Role) {
					return appErrors.Unauthorized("")
				}

				ctx := context.WithValue(req.Context(), "user", &userData)
				req = req.WithContext(ctx)
			}

			return next(res, req)
		}
	}
}
