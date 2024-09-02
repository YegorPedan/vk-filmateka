package mockRepository

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
)

type filmRepository struct {
	db *inMemDb.InMemDb
}

func (f filmRepository) Create(ctx context.Context, aggregate *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error) {
	has := slices.ContainsFunc(f.db.Film, func(item *model.Film) bool {
		if aggregate.Film.Id == item.Id {
			return true
		}
		return false
	})

	if has {
		return nil, errors.New("conflict fields")
	}

	film := model.Film{Id: aggregate.Film.Id, Name: aggregate.Film.Name, ReleaseDate: aggregate.Film.ReleaseDate, Rate: aggregate.Film.Rate, Description: aggregate.Film.Description}
	f.db.Film = append(f.db.Film, &film)
	return aggregate, nil
}

func (f filmRepository) Update(ctx context.Context, aggregate *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error) {
	has := slices.ContainsFunc(f.db.Film, func(item *model.Film) bool {
		if aggregate.Film.Id == item.Id {
			return true
		}
		return false
	})

	if !has {
		return nil, sql.ErrNoRows
	}

	for i, item := range f.db.Film {
		if aggregate.Film.Id == item.Id {
			film := aggregate.Film
			f.db.Film[i] = &film
		}
	}

	return aggregate, nil
}

func (f filmRepository) Delete(ctx context.Context, id string) error {
	has := false

	var filteredFilms []*model.Film

	for _, film := range f.db.Film {
		if film.Id != id {
			filteredFilms = append(filteredFilms, film)
		} else {
			has = true
		}
	}

	f.db.Film = filteredFilms

	if has {
		return nil
	}
	return sql.ErrNoRows
}

func (f filmRepository) GetById(ctx context.Context, id string) (*aggregate.FilmAggregate, error) {
	var searched *model.Film
	for _, film := range f.db.Film {
		if film.Id == id {
			searched = film
		}
	}
	if searched != nil {
		return &aggregate.FilmAggregate{Film: *searched}, nil
	}
	return nil, sql.ErrNoRows
}

// TODO - Нету обработки на сортировку по полю
func (f filmRepository) GetByQuery(ctx context.Context, query domainQuery.FilmRepositoryQuery) ([]*aggregate.FilmAggregate, int, error) {
	if len(f.db.Film) == 0 {
		return []*aggregate.FilmAggregate{}, 0, nil
	}

	start := query.PageCount * (query.CurrentPage - 1)
	if start >= len(f.db.Film) {
		return nil, 0, errors.New("current page incorrect")
	}

	getted := make([]*aggregate.FilmAggregate, 0, query.PageCount)

	for i, j := 0, start; j < len(f.db.Film) && i < query.PageCount; i++ {
		var actors []*model.Actor = nil
		isActorConnection := slices.ContainsFunc(query.WithConnection, func(item string) bool {
			if item == "actor" {
				return true
			}
			return false
		})
		if isActorConnection {
			actorIds := make([]string, 0, 50)
			for _, item := range f.db.ActorFilm {
				if item.FilmId == f.db.Film[j].Id {
					actorIds = append(actorIds, item.ActorId)
				}
			}

			for _, actor := range f.db.Actor {
				has := slices.ContainsFunc(actorIds, func(id string) bool {
					if id == actor.Id {
						return true
					}
					return false
				})
				if has {
					cpy := *actor
					actors = append(actors, &cpy)
				}
			}
		}
		aggr := aggregate.FilmAggregate{
			Film:   *f.db.Film[j],
			Actors: actors,
		}
		getted = append(getted, &aggr)
		j++
	}

	remainder := len(f.db.Film) % query.PageCount
	totalPageCount := len(f.db.Film) / query.PageCount
	if remainder != 0 {
		totalPageCount++
	}

	return getted, totalPageCount, nil
}

func (f filmRepository) SearchByNameAndActorName(ctx context.Context, searchValue string) ([]*aggregate.FilmAggregate, int, error) {
	foundItems := make([]*aggregate.FilmAggregate, 0, 100)

	for _, film := range f.db.Film {
		if strings.Contains(film.Name, searchValue) {
			aggr := aggregate.FilmAggregate{Film: *film}
			foundItems = append(foundItems, &aggr)
		}
	}

	remainder := len(f.db.Film) % 15
	totalPageCount := len(f.db.Film) / 15
	if remainder != 0 {
		totalPageCount++
	}

	return foundItems, totalPageCount, nil
}

func NewFilmRepository() repository.FilmRepository {
	return &filmRepository{inMemDb.New()}
}
