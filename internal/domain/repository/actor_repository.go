package repository

import (
	"context"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
)

type ActorRepository interface {
	Create(ctx context.Context, aggregate *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error)
	Update(ctx context.Context, aggregate *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error)
	Delete(ctx context.Context, id string) error
	AddFilm(ctx context.Context, actorId string, filmIds ...string) error
	GetById(ctx context.Context, id string) (*aggregate.ActorAggregate, error)
	GetByQuery(ctx context.Context, query domainQuery.ActorRepositoryQuery) ([]*aggregate.ActorAggregate, int, error)
}
