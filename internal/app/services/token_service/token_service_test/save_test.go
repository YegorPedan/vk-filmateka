package token_servicetest_test

import (
	"context"
	"testing"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTokenServiceSave(t *testing.T) {
	config.MustLoad()
	tokenRepo := mockRepository.NewTokenRepository()
	tokenServ := tokenService.New(tokenRepo)

	id := uuid.New().String()
	jwtData := tokenService.JwtUserData{Id: id, Role: constants.UserRole}
	tokens, err := tokenServ.Generate(jwtData)
	if err != nil {
		t.Fatal(err)
	}
	has, _ := tokenRepo.HasByValue(context.Background(), tokens.RefreshToken)

	assert.False(t, has)
	tokenServ.Save(context.Background(), appDto.SaveTokenServiceDto{Id: jwtData.Id, RefreshToken: tokens.RefreshToken})
	has, _ = tokenRepo.HasByValue(context.Background(), tokens.RefreshToken)
	assert.True(t, has)

	jwtData.Role = constants.AdminRole

	tokens, err = tokenServ.Generate(jwtData)
	if err != nil {
		t.Fatal(err)
	}
	has, _ = tokenRepo.HasByValue(context.Background(), tokens.RefreshToken)
	assert.False(t, has)
	_, err = tokenServ.Save(context.Background(), appDto.SaveTokenServiceDto{Id: jwtData.Id, RefreshToken: tokens.RefreshToken})
	if err != nil {
		t.Fatal(err)
	}
	has, _ = tokenRepo.HasByValue(context.Background(), tokens.RefreshToken)
	assert.True(t, has)
}
