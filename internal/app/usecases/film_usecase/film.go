package filmUseCase

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
	FilmUseCase interface {
		Create(ctx context.Context, data appDto.CreateFilmUseCaseDto) (*aggregate.FilmAggregate, error)
		Update(ctx context.Context, data *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error)
		Delete(ctx context.Context, id string) error
		GetById(ctx context.Context, id string) (*aggregate.FilmAggregate, error)
		GetByQuery(ctx context.Context, query domainQuery.FilmRepositoryQuery) (*appDto.FilmGetByQueryResult, error)
		SearchByNameAndActorName(ctx context.Context, searchValue string) (*appDto.FilmGetByQueryResult, error)
	}

	filmUseCase struct {
		repository.FilmRepository
	}
)

func (f filmUseCase) Create(ctx context.Context, data appDto.CreateFilmUseCaseDto) (*aggregate.FilmAggregate, error) {
	filmAggregate, err := aggregate.NewFilmAggregate(model.Film{
		Id:          uuid.New().String(),
		Name:        data.Name,
		Description: data.Description,
		ReleaseDate: data.ReleaseDate,
		Rate:        data.Rate,
	})

	if err != nil {
		return nil, appErrors.UnprocessableEntity("", "target: ActorUseCase, method: Create ", "error: ", err.Error())
	}

	createAggregate, err := f.FilmRepository.Create(ctx, filmAggregate)
	if err != nil {
		return nil, appErrors.InternalServerError("", "target: ActorUseCase, method: Create ", "repository create error: ", err.Error())
	}

	return createAggregate, nil
}

func (f filmUseCase) Update(ctx context.Context, data *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error) {
	user, _ := f.FilmRepository.GetById(ctx, data.Film.Id)
	if user == nil {
		return nil, appErrors.NotFound("")
	}

	updateAggregate, err := f.FilmRepository.Update(ctx, data)
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}

	return updateAggregate, err
}

func (f filmUseCase) Delete(ctx context.Context, id string) error {
	err := f.FilmRepository.Delete(ctx, id)
	if err == sql.ErrNoRows {
		return appErrors.NotFound("", "error: ", err.Error())
	}
	if err != nil {
		return appErrors.InternalServerError("", "error: ", err.Error())
	}

	return nil
}

func (f filmUseCase) GetById(ctx context.Context, id string) (*aggregate.FilmAggregate, error) {
	film, err := f.FilmRepository.GetById(ctx, id)
	if err == sql.ErrNoRows {
		return nil, appErrors.NotFound("")
	}
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}
	return film, nil
}

func (f filmUseCase) GetByQuery(ctx context.Context, query domainQuery.FilmRepositoryQuery) (*appDto.FilmGetByQueryResult, error) {
	byQuery, pageCount, err := f.FilmRepository.GetByQuery(ctx, query)
	if err != nil {
		return nil, appErrors.InternalServerError("", "error:", err.Error())
	}

	return &appDto.FilmGetByQueryResult{
		Films:     byQuery,
		PageCount: pageCount,
	}, nil
}

func (f filmUseCase) SearchByNameAndActorName(ctx context.Context, searchValue string) (*appDto.FilmGetByQueryResult, error) {
	films, pageCount, err := f.FilmRepository.SearchByNameAndActorName(ctx, searchValue)
	if len(films) == 0 {
		return nil, appErrors.NotFound("Search film not found")
	}
	if err != nil {
		return nil, appErrors.InternalServerError("")
	}
	return &appDto.FilmGetByQueryResult{
		Films:     films,
		PageCount: pageCount,
	}, nil
}

func New(filmRepository repository.FilmRepository) FilmUseCase {
	return &filmUseCase{
		FilmRepository: filmRepository,
	}
}
