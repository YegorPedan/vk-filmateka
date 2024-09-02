package authUseCase

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	appMapper "github.com/OddEer0/vk-filmoteka/internal/app/app_mapper"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
)

func (a *authUseCase) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	if refreshToken == "" {
		return nil, appErrors.Unauthorized(constants.Unauthorized)
	}

	jwtUserData, err := a.TokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	has, err := a.TokenService.HasByValue(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, appErrors.Unauthorized(constants.Unauthorized)
	}
	userAggregate, err := a.UserRepository.GetById(ctx, jwtUserData.Id)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: AuthUseCase, method: Refresh", "get user by id error", err.Error())
	}

	tokens, err := a.TokenService.Generate(*jwtUserData)
	if err != nil {
		return nil, err
	}
	_, err = a.TokenService.Save(ctx, appDto.SaveTokenServiceDto{Id: userAggregate.User.Id, RefreshToken: tokens.RefreshToken})
	if err != nil {
		return nil, err
	}
	userMapper := appMapper.NewUserAggregateMapper()
	responseUser := userMapper.ToResponseUserDto(userAggregate)

	return &AuthResult{User: &responseUser, Tokens: *tokens}, nil
}
