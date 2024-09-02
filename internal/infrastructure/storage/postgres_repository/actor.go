package postgresRepository

import (
	"context"
	"database/sql"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	"slices"
	"time"
)

type actorRepository struct {
	db *sql.DB
}

func (a actorRepository) Create(ctx context.Context, data *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error) {
	query := "INSERT INTO actors (id, name, gender, birthday) VALUES ($1, $2, $3, $4) RETURNING id, name, gender, birthday"
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	actor := data.Actor
	err = stmt.QueryRowContext(ctx, actor.Id, actor.Name, actor.Gender, actor.Birthday).Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.Birthday)
	if err != nil {
		return nil, err
	}

	data.Actor = actor
	return data, nil
}

func (a actorRepository) Update(ctx context.Context, data *aggregate.ActorAggregate) (*aggregate.ActorAggregate, error) {
	query := "UPDATE actors SET name = $1, gender = $2, birthday = $3 WHERE id = $4 RETURNING id, name, gender, birthday"
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	actor := data.Actor
	err = stmt.QueryRowContext(ctx, actor.Name, actor.Gender, actor.Birthday, actor.Id).Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.Birthday)
	if err != nil {
		return nil, err
	}

	data.Actor = actor
	return data, nil
}

func (a actorRepository) Delete(ctx context.Context, id string) error {
	_, err := a.db.ExecContext(ctx, "DELETE FROM actor_film WHERE actor_id = $1", id)
	if err != nil {
		return err
	}
	query := "DELETE FROM actors WHERE id = $1"
	_, err = a.db.ExecContext(ctx, query, id)
	return err
}

func (a actorRepository) AddFilm(ctx context.Context, actorId string, filmIds ...string) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var actorExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM actors WHERE id = $1)", actorId).Scan(&actorExists)
	if err != nil {
		return err
	}
	if !actorExists {
		return sql.ErrNoRows
	}

	for _, filmId := range filmIds {
		var filmExists bool
		err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM films WHERE id = $1)", filmId).Scan(&filmExists)
		if err != nil {
			return err
		}
		if !filmExists {
			return sql.ErrNoRows
		}

		_, err := tx.ExecContext(ctx, "INSERT INTO actor_film (actor_id, film_id) VALUES ($1, $2)", actorId, filmId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a actorRepository) GetById(ctx context.Context, id string) (*aggregate.ActorAggregate, error) {
	query := "SELECT id, name, gender, birthday FROM actors WHERE id = $1"
	row := a.db.QueryRowContext(ctx, query, id)

	var actor model.Actor
	err := row.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.Birthday)
	if err != nil {
		return nil, err
	}

	return &aggregate.ActorAggregate{Actor: actor}, nil
}

func (a actorRepository) GetByQuery(ctx context.Context, query domainQuery.ActorRepositoryQuery) ([]*aggregate.ActorAggregate, int, error) {
	offset := query.PageCount * (query.CurrentPage - 1)
	limit := query.PageCount

	sqlQuery := `
		SELECT
		a.*, 
		f.id AS film_id, 
		f.name AS film_name, 
		f.description AS film_description, 
		f.release_date AS film_release_date, 
		f.rate AS film_rate
		FROM actors a
		LEFT JOIN actor_film af ON a.id = af.actor_id and 'film' = $1
		LEFT JOIN films f ON af.film_id = f.id and 'film' = $1
		ORDER BY a.id, f.id NULLS LAST
		LIMIT $2 OFFSET $3
    `

	conn := "not"
	if slices.Contains(query.WithConnection, "film") {
		conn = "film"
	}

	rows, err := a.db.QueryContext(ctx, sqlQuery, conn, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	actorMap := make(map[string]*aggregate.ActorAggregate, 20)
	actors := make([]*aggregate.ActorAggregate, 0, len(actorMap))

	for rows.Next() {
		var (
			actorID         string
			actorName       string
			actorGender     string
			actorBirthday   time.Time
			filmID          sql.NullString
			filmName        sql.NullString
			filmDescription sql.NullString
			filmReleaseDate sql.NullTime
			filmRate        sql.NullFloat64
		)

		err := rows.Scan(&actorID, &actorName, &actorGender, &actorBirthday, &filmID, &filmName, &filmDescription, &filmReleaseDate, &filmRate)
		if err != nil {
			return nil, 0, err
		}

		var aggr *aggregate.ActorAggregate
		if res, ok := actorMap[actorID]; !ok {
			aggr = &aggregate.ActorAggregate{
				Actor: model.Actor{
					Id:       actorID,
					Name:     actorName,
					Gender:   actorGender,
					Birthday: actorBirthday,
				},
			}
			actorMap[actorID] = aggr
			actors = append(actors, aggr)
		} else {
			aggr = res
		}

		if filmID.Valid {
			aggr.Films = append(aggr.Films, &model.Film{Id: filmID.String, Name: filmName.String, Description: &filmDescription.String, ReleaseDate: filmReleaseDate.Time, Rate: float32(filmRate.Float64)})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	totalCount := 0
	err = a.db.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM actors
    `).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	remainder := totalCount % limit
	totalCount /= limit
	if remainder != 0 {
		totalCount++
	}

	return actors, totalCount, nil
}

func NewActorRepository(db *sql.DB) repository.ActorRepository {
	return &actorRepository{db: db}
}
