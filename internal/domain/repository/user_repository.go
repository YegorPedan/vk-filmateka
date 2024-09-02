package repository

import (
	"context"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
)

type UserRepository interface {
	Create(ctx context.Context, userAggregate *aggregate.UserAggregate) (*aggregate.UserAggregate, error)
	Update(ctx context.Context, userAggregate *aggregate.UserAggregate) (*aggregate.UserAggregate, error)
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (*aggregate.UserAggregate, error)
	HasUserByName(ctx context.Context, name string) (bool, error)
	GetByName(ctx context.Context, name string) (*aggregate.UserAggregate, error)
}
