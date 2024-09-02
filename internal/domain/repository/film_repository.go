package repository

import (
	"context"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
)

type FilmRepository interface {
	Create(ctx context.Context, aggregate *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error)
	Update(ctx context.Context, aggregate *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error)
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (*aggregate.FilmAggregate, error)
	GetByQuery(ctx context.Context, query domainQuery.FilmRepositoryQuery) ([]*aggregate.FilmAggregate, int, error)
	SearchByNameAndActorName(ctx context.Context, searchValue string) ([]*aggregate.FilmAggregate, int, error)
}
