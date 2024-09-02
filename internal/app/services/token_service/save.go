package tokenService

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
)

func (t *tokenService) Save(ctx context.Context, data appDto.SaveTokenServiceDto) (*model.Token, error) {
	_, err := t.TokenRepository.GetById(ctx, data.Id)

	if err != nil {
		token, err := t.TokenRepository.Create(ctx, &model.Token{Id: data.Id, Value: data.RefreshToken})
		if err != nil {
			return nil, appErrors.InternalServerError("")
		}
		return token, nil
	}

	token, err := t.TokenRepository.Update(ctx, &model.Token{Id: data.Id, Value: data.RefreshToken})
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}
	return token, nil
}
