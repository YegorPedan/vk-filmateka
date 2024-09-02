package userService

import (
	"context"

	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/valuesobject"
	"github.com/google/uuid"
)

// Create method create user aggregate, does not create db table
func (u *userService) Create(ctx context.Context, data appDto.RegistrationUseCaseDto) (*aggregate.UserAggregate, error) {
	candidate, err := u.UserRepository.HasUserByName(ctx, data.Name)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: UserService, method: Create. ", "user repository error: ", err.Error())
	}
	if candidate {
		return nil, appErrors.Conflict(constants.UserNickExist, "target: UserService, method: Create. ", "Nick conflict")
	}

	hashPassword, err := valuesobject.NewPassword(data.Password)
	if err != nil {
		return nil, appErrors.UnprocessableEntity("", "target: UserService, method: Create. ", "valuesobject NewPassword method error: ", err.Error())
	}
	userAggregate, err := aggregate.NewUserAggregate(model.User{
		Id:       uuid.New().String(),
		Name:     data.Name,
		Role:     constants.UserRole,
		Password: hashPassword,
	})
	if err != nil {
		return nil, appErrors.UnprocessableEntity("", "target: UserService, method: Create. ", "aggregate CreateNewAggregate error: ", err.Error())
	}

	return userAggregate, nil
}
