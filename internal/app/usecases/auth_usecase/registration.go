package authUseCase

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	appMapper "github.com/OddEer0/vk-filmoteka/internal/app/app_mapper"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
)

func (a *authUseCase) Registration(ctx context.Context, data appDto.RegistrationUseCaseDto) (*AuthResult, error) {
	userAggregate, err := a.UserService.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	tokens, err := a.TokenService.Generate(tokenService.JwtUserData{Id: userAggregate.User.Id, Role: userAggregate.User.Role})
	if err != nil {
		return nil, err
	}
	err = userAggregate.SetToken(tokens.RefreshToken)
	if err != nil {
		return nil, appErrors.UnprocessableEntity("", "target: AuthUseCase, method: Registration. ", "Aggregate SetToken method error: ", err.Error())
	}

	dbUserAggregate, err := a.UserRepository.Create(ctx, userAggregate)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: AuthUseCase, method: Registration. ", "UserRepository create user error: ", err.Error())
	}
	userMapper := appMapper.NewUserAggregateMapper()
	responseUser := userMapper.ToResponseUserDto(dbUserAggregate)

	return &AuthResult{User: &responseUser, Tokens: *tokens}, nil
}
