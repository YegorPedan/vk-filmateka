package postgresRepository

import (
	"context"
	"database/sql"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	"slices"
)

type filmRepository struct {
	db *sql.DB
}

func (f filmRepository) Create(ctx context.Context, aggregate *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error) {
	query := "INSERT INTO films (id, name, description, release_date, rate) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, release_date, rate"
	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	film := aggregate.Film
	err = stmt.QueryRowContext(ctx, film.Id, film.Name, film.Description, film.ReleaseDate, film.Rate).Scan(&film.Id, &film.Name, &film.Description, &film.ReleaseDate, &film.Rate)
	if err != nil {
		return nil, err
	}

	aggregate.Film = film
	return aggregate, nil
}

func (f filmRepository) Update(ctx context.Context, aggregate *aggregate.FilmAggregate) (*aggregate.FilmAggregate, error) {
	query := "UPDATE films SET name = $1, description = $2, release_date = $3, rate = $4 WHERE id = $5 RETURNING id, name, description, release_date, rate"
	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	film := aggregate.Film
	err = stmt.QueryRowContext(ctx, film.Name, film.Description, film.ReleaseDate, film.Rate, film.Id).Scan(&film.Id, &film.Name, &film.Description, &film.ReleaseDate, &film.Rate)
	if err != nil {
		return nil, err
	}

	aggregate.Film = film
	return aggregate, nil
}

func (f filmRepository) Delete(ctx context.Context, id string) error {
	_, err := f.db.ExecContext(ctx, "DELETE FROM actor_film WHERE film_id = $1", id)
	if err != nil {
		return err
	}
	query := "DELETE FROM films WHERE id = $1"
	_, err = f.db.ExecContext(ctx, query, id)
	return err
}

func (f filmRepository) GetById(ctx context.Context, id string) (*aggregate.FilmAggregate, error) {
	query := "SELECT id, name, description, release_date, rate FROM films WHERE id = $1"
	row := f.db.QueryRowContext(ctx, query, id)

	var film model.Film
	err := row.Scan(&film.Id, &film.Name, &film.Description, &film.ReleaseDate, &film.Rate)
	if err != nil {
		return nil, err
	}

	return &aggregate.FilmAggregate{Film: film}, nil
}

func (f filmRepository) GetByQuery(ctx context.Context, query domainQuery.FilmRepositoryQuery) ([]*aggregate.FilmAggregate, int, error) {
	offset := query.PageCount * (query.CurrentPage - 1)
	limit := query.PageCount
	field := "f.rate"
	if query.SortField == "name" {
		field = "f.name"
	} else if query.SortField == "released_date" {
		field = "f.released_date"
	}

	sqlQuery := `
        SELECT
            f.*,
            a.id AS actor_id,
            a.name AS actor_name,
            a.gender AS actor_gender,
            a.birthday AS actor_birthday
        FROM films f
        LEFT JOIN actor_film af ON f.id = af.film_id AND 'actor' = $1
        LEFT JOIN actors a ON af.actor_id = a.id AND 'actor' = $1
		ORDER BY ` + field + " " + string(query.OrderBy) + `
		LIMIT $2 OFFSET $3;
		`

	conn := "not"
	if slices.Contains(query.WithConnection, "actor") {
		conn = "actor"
	}

	rows, err := f.db.QueryContext(ctx, sqlQuery, conn, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	filmsMap := make(map[string]*aggregate.FilmAggregate, 20)
	films := make([]*aggregate.FilmAggregate, 0, len(filmsMap))

	for rows.Next() {
		var (
			film          model.Film
			actorId       sql.NullString
			actorName     sql.NullString
			actorBirthday sql.NullTime
			actorGender   sql.NullString
		)

		err := rows.Scan(
			&film.Id,
			&film.Name,
			&film.Description,
			&film.ReleaseDate,
			&film.Rate,
			&actorId,
			&actorName,
			&actorGender,
			&actorBirthday,
		)

		if err != nil {
			return nil, 0, err
		}

		var aggr *aggregate.FilmAggregate
		if res, ok := filmsMap[film.Id]; !ok {
			aggr = &aggregate.FilmAggregate{
				Film: film,
			}
			filmsMap[film.Id] = aggr
			films = append(films, aggr)
		} else {
			aggr = res
		}

		if actorId.Valid {
			aggr.Actors = append(aggr.Actors, &model.Actor{Id: actorId.String, Name: actorName.String, Birthday: actorBirthday.Time, Gender: actorGender.String})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	totalCount := 0
	err = f.db.QueryRowContext(ctx, `
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

	return films, totalCount, nil
}

func (f filmRepository) SearchByNameAndActorName(ctx context.Context, searchValue string) ([]*aggregate.FilmAggregate, int, error) {
	sqlQuery := `
		SELECT DISTINCT ON(films.id) films.*
		FROM films
		LEFT JOIN actor_film af ON films.id = af.film_id
		LEFT JOIN actors a ON af.actor_id = a.id
		WHERE (films.name LIKE '%' || $1 || '%' OR a.name LIKE '%' || $1 || '%')
		ORDER BY films.id, films.rate ASC
		LIMIT 20 OFFSET 0;
	`

	rows, err := f.db.QueryContext(ctx, sqlQuery, searchValue)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	films := make([]*aggregate.FilmAggregate, 0, 20)

	for rows.Next() {
		aggr := &aggregate.FilmAggregate{}
		err := rows.Scan(&aggr.Film.Id, &aggr.Film.Name, &aggr.Film.Description, &aggr.Film.ReleaseDate, &aggr.Film.Rate)
		if err != nil {
			return nil, 0, err
		}

		films = append(films, aggr)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return films, 1, nil
}

func NewFilmRepository(db *sql.DB) repository.FilmRepository {
	return &filmRepository{db: db}
}
