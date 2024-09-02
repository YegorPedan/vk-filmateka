package authUseCase

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	appMapper "github.com/OddEer0/vk-filmoteka/internal/app/app_mapper"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"golang.org/x/crypto/bcrypt"
)

func (a *authUseCase) Login(ctx context.Context, data appDto.LoginUseCaseDto) (*AuthResult, error) {
	candidate, err := a.UserRepository.HasUserByName(ctx, data.Name)
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}
	if !candidate {
		return nil, appErrors.Forbidden(constants.NickOrPasswordIncorrect)
	}
	userAggregate, err := a.UserRepository.GetByName(ctx, data.Name)
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}

	isEqualPassword := bcrypt.CompareHashAndPassword([]byte(userAggregate.User.Password.Value), []byte(data.Password))
	if isEqualPassword != nil {
		return nil, appErrors.Forbidden(constants.NickOrPasswordIncorrect)
	}

	tokens, err := a.TokenService.Generate(tokenService.JwtUserData{Id: userAggregate.User.Id, Role: userAggregate.User.Role})
	if err != nil {
		return nil, err
	}

	_, err = a.TokenService.Save(ctx, appDto.SaveTokenServiceDto{Id: userAggregate.User.Id, RefreshToken: tokens.RefreshToken})
	if err != nil {
		return nil, err
	}

	aggregateMapper := appMapper.NewUserAggregateMapper()
	responseUser := aggregateMapper.ToResponseUserDto(userAggregate)

	return &AuthResult{User: &responseUser, Tokens: *tokens}, nil
}
