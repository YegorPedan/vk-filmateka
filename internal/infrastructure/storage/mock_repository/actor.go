package mockRepository

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
)

type actorRepository struct {
	db *inMemDb.InMemDb
}

func (a actorRepository) Create(ctx context.Context, aggregate *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error) {
	has := slices.ContainsFunc(a.db.Actor, func(item *model.Actor) bool {
		if aggregate.Actor.Id == item.Id {
			return true
		}
		return false
	})

	if has {
		return nil, errors.New("conflict fields")
	}

	actor := model.Actor{Id: aggregate.Actor.Id, Name: aggregate.Actor.Name, Gender: aggregate.Actor.Gender, Birthday: aggregate.Actor.Birthday}
	a.db.Actor = append(a.db.Actor, &actor)
	return aggregate, nil
}

func (a actorRepository) Update(ctx context.Context, aggregate *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error) {
	has := slices.ContainsFunc(a.db.Actor, func(item *model.Actor) bool {
		if aggregate.Actor.Id == item.Id {
			return true
		}
		return false
	})

	if !has {
		return nil, sql.ErrNoRows
	}

	for i, item := range a.db.Actor {
		if aggregate.Actor.Id == item.Id {
			actor := aggregate.Actor
			a.db.Actor[i] = &actor
		}
	}

	return aggregate, nil
}

func (a actorRepository) Delete(ctx context.Context, id string) error {
	has := false

	var filteredActors []*model.Actor

	for _, actor := range a.db.Actor {
		if actor.Id != id {
			filteredActors = append(filteredActors, actor)
		} else {
			has = true
		}
	}

	a.db.Actor = filteredActors

	if has {
		return nil
	}
	return sql.ErrNoRows
}

func (a actorRepository) AddFilm(ctx context.Context, actorId string, filmIds ...string) error {
	added := make([]*inMemDb.ActorFilm, 0, len(filmIds))
	for _, id := range filmIds {
		var searchedFilm *model.Film
		for _, film := range a.db.Film {
			if film.Id == id {
				searchedFilm = film
			}
		}

		if searchedFilm == nil {
			return sql.ErrNoRows
		}

		var searchedActor *model.Actor
		for _, actor := range a.db.Actor {
			if actor.Id == actorId {
				searchedActor = actor
			}
		}

		if searchedActor == nil {
			return sql.ErrNoRows
		}

		added = append(added, &inMemDb.ActorFilm{
			ActorId: searchedActor.Id,
			FilmId:  searchedFilm.Id,
		})
	}

	a.db.ActorFilm = append(a.db.ActorFilm, added...)
	return nil
}

func (a actorRepository) GetById(ctx context.Context, id string) (*aggregate.ActorAggregate, error) {
	var searched *model.Actor
	for _, actor := range a.db.Actor {
		if actor.Id == id {
			searched = actor
		}
	}
	if searched != nil {
		return &aggregate.ActorAggregate{Actor: *searched}, nil
	}
	return nil, sql.ErrNoRows
}

func (a actorRepository) GetByQuery(ctx context.Context, query domainQuery.ActorRepositoryQuery) ([]*aggregate.ActorAggregate, int, error) {
	if len(a.db.Actor) == 0 {
		return []*aggregate.ActorAggregate{}, 0, nil
	}

	start := query.PageCount * (query.CurrentPage - 1)
	if start >= len(a.db.Actor) {
		return nil, 0, errors.New("current page incorrect")
	}

	getted := make([]*aggregate.ActorAggregate, 0, query.PageCount)

	for i, j := 0, start; j < len(a.db.Actor) && i < query.PageCount; i++ {
		var films []*model.Film = nil
		isFilmConnection := slices.ContainsFunc(query.WithConnection, func(item string) bool {
			if item == "film" {
				return true
			}
			return false
		})
		if isFilmConnection {
			filmIds := make([]string, 0, 50)
			for _, item := range a.db.ActorFilm {
				if item.ActorId == a.db.Actor[j].Id {
					filmIds = append(filmIds, item.FilmId)
				}
			}

			for _, film := range a.db.Film {
				has := slices.ContainsFunc(filmIds, func(id string) bool {
					if id == film.Id {
						return true
					}
					return false
				})
				if has {
					cpy := *film
					films = append(films, &cpy)
				}
			}
		}
		aggr := aggregate.ActorAggregate{
			Actor: *a.db.Actor[j],
			Films: films,
		}
		getted = append(getted, &aggr)
		j++
	}

	remainder := len(a.db.Actor) % query.PageCount
	totalPageCount := len(a.db.Actor) / query.PageCount
	if remainder != 0 {
		totalPageCount++
	}

	return getted, totalPageCount, nil
}

func NewActorRepository() repository.ActorRepository {
	return &actorRepository{inMemDb.New()}
}
