package tokenService

import (
	"context"

	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
)

func (t *tokenService) HasByValue(ctx context.Context, refreshToken string) (bool, error) {
	result, err := t.TokenRepository.HasByValue(ctx, refreshToken)
	if err != nil {
		return result, appErrors.InternalServerError("", "target: TokenService, method: HasByValue. ", "error: ", err.Error())
	}
	return result, nil
}
