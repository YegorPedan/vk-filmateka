package appMapper

import (
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
)

type (
	UserAggregateMapper interface {
		ToResponseUserDto(aggregate *aggregate.UserAggregate) appDto.ResponseUserDto
	}

	userAggregateMapper struct{}
)

func NewUserAggregateMapper() UserAggregateMapper {
	return userAggregateMapper{}
}

func (u userAggregateMapper) ToResponseUserDto(aggregate *aggregate.UserAggregate) appDto.ResponseUserDto {
	return appDto.ResponseUserDto{
		Id:   aggregate.User.Id,
		Name: aggregate.User.Name,
		Role: aggregate.User.Role,
	}
}
