package actorUseCase

import (
	"context"
	"database/sql"
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	"github.com/google/uuid"
)

type (
	ActorUseCase interface {
		Create(ctx context.Context, data appDto.CreateActorUseCaseDto) (*aggregate.ActorAggregate, error)
		Update(ctx context.Context, data *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error)
		Delete(ctx context.Context, id string) error
		GetById(ctx context.Context, id string) (*aggregate.ActorAggregate, error)
		GetByQuery(ctx context.Context, query domainQuery.ActorRepositoryQuery) (*appDto.ActorGetByQueryResult, error)
		AddFilm(ctx context.Context, actorId string, filmIds ...string) error
	}

	actorUseCase struct {
		repository.ActorRepository
		repository.FilmRepository
	}
)

func (a *actorUseCase) Create(ctx context.Context, data appDto.CreateActorUseCaseDto) (*aggregate.ActorAggregate, error) {
	actorAggregate, err := aggregate.NewActorAggregate(model.Actor{
		Id:       uuid.New().String(),
		Name:     data.Name,
		Gender:   data.Gender,
		Birthday: data.Birthday,
	})

	if err != nil {
		return nil, appErrors.UnprocessableEntity("", "target: ActorUseCase, method: Create ", "error: ", err.Error())
	}

	createAggregate, err := a.ActorRepository.Create(ctx, actorAggregate)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: ActorUseCase, method: Create ", "repository create error: ", err.Error())
	}

	return createAggregate, nil
}

func (a *actorUseCase) Update(ctx context.Context, data *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error) {
	updateAggregate, err := a.ActorRepository.Update(ctx, data)
	if err == sql.ErrNoRows {
		return nil, appErrors.NotFound("")
	}
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}

	return updateAggregate, err
}

func (a *actorUseCase) Delete(ctx context.Context, id string) error {
	err := a.ActorRepository.Delete(ctx, id)
	if err == sql.ErrNoRows {
		return appErrors.NotFound("")
	} else if err != nil {
		return appErrors.InternalServerError("")
	}
	return nil
}

func (a *actorUseCase) GetByQuery(ctx context.Context, query domainQuery.ActorRepositoryQuery) (*appDto.ActorGetByQueryResult, error) {
	byQuery, pageCount, err := a.ActorRepository.GetByQuery(ctx, query)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: ActorUseCase, method: GetByQuery. ", "get by query error: ", err.Error())
	}
	return &appDto.ActorGetByQueryResult{
		Actors:    byQuery,
		PageCount: pageCount,
	}, nil
}

func (a *actorUseCase) GetById(ctx context.Context, id string) (*aggregate.ActorAggregate, error) {
	byId, err := a.ActorRepository.GetById(ctx, id)
	if err == sql.ErrNoRows {
		return nil, appErrors.NotFound("")
	}
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}
	return byId, nil
}

func (a *actorUseCase) AddFilm(ctx context.Context, actorId string, filmIds ...string) error {
	if len(filmIds) == 0 {
		return appErrors.InternalServerError("", "target: ActorUseCase, method: AddFilm", "error: ", "not id or ids")
	}
	err := a.ActorRepository.AddFilm(ctx, actorId, filmIds...)
	if err == sql.ErrNoRows {
		return appErrors.NotFound("")
	}
	if err != nil {
		return appErrors.InternalServerError("", "target: ActorUseCase, method: AddFilm", " added repository error: ", err.Error())
	}

	return nil
}

func New(actorRepository repository.ActorRepository, filmRepository repository.FilmRepository) ActorUseCase {
	return &actorUseCase{
		ActorRepository: actorRepository,
		FilmRepository:  filmRepository,
	}
}
