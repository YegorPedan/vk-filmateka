package auth_usecase_test

import (
	"context"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	userService "github.com/OddEer0/vk-filmoteka/internal/app/services/user_service"
	authUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/auth_usecase"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuthRefresh(t *testing.T) {
	mockData := newAuthUseCaseDataMock()
	userRepo := mockRepository.NewUserRepository()
	tokenRepo := mockRepository.NewTokenRepository()

	useCase := authUseCase.New(userService.New(userRepo), tokenService.New(tokenRepo), userRepo)
	config.MustLoad()
	_, _ = useCase.Registration(context.Background(), mockData.Registration.CorrectRegInput1)
	res, _ := useCase.Login(context.Background(), mockData.Login.CorrectRegInput1)
	refresh, err := useCase.Refresh(context.Background(), res.Tokens.RefreshToken)
	assert.Nil(t, err)
	assert.Equal(t, refresh.User, res.User)
	if err != nil {
		t.Fatal("error refresh")
	}

	err = useCase.Logout(context.Background(), res.Tokens.RefreshToken)
	assert.Nil(t, err)
	refresh, err = useCase.Refresh(context.Background(), res.Tokens.RefreshToken)
	assert.NotNil(t, err)
	assert.Nil(t, refresh)
	res, _ = useCase.Login(context.Background(), mockData.Login.CorrectRegInput1)
	time.Sleep(6 * time.Second)
	refresh, err = useCase.Refresh(context.Background(), res.Tokens.RefreshToken)
	assert.NotNil(t, err)
	assert.Nil(t, refresh)
	refresh, err = useCase.Refresh(context.Background(), "")
	assert.NotNil(t, err)
	assert.Nil(t, refresh)

	db := inMemDb.New()
	db.CleanUp()
}
