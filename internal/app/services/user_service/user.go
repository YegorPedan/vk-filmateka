package userService

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
)

type (
	Service interface {
		Create(ctx context.Context, data appDto.RegistrationUseCaseDto) (*aggregate.UserAggregate, error)
	}

	userService struct {
		repository.UserRepository
	}
)

func New(userRepo repository.UserRepository) Service {
	return &userService{UserRepository: userRepo}
}
