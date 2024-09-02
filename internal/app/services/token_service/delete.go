package tokenService

import (
	"context"

	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
)

func (t *tokenService) DeleteByValue(ctx context.Context, value string) error {
	err := t.TokenRepository.DeleteByValue(ctx, value)
	if err != nil {
		return appErrors.InternalServerError("", "target: TokenService, method: DeleteByValue", "error: ", err.Error())
	}
	return nil
}
