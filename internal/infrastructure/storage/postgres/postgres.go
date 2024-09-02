package postgres

import (
	"database/sql"
	"fmt"
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	"github.com/OddEer0/vk-filmoteka/internal/domain/valuesobject"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func ConnectPg(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS actors (
        id UUID PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        gender VARCHAR(10) NOT NULL,
        birthday DATE NOT NULL
    )`); err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS films (
        id UUID PRIMARY KEY,
        name VARCHAR(150) NOT NULL,
        description TEXT,
        release_date DATE NOT NULL,
        rate NUMERIC(4,1) NOT NULL
    )`); err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS actor_film (
        actor_id UUID REFERENCES actors(id),
		film_id UUID REFERENCES films(id),
		PRIMARY KEY (actor_id, film_id)
    )`); err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(20) NOT NULL
    )`); err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
        id UUID PRIMARY KEY REFERENCES users(id),
    	value VARCHAR(255) NOT NULL
    )`); err != nil {
		return nil, err
	}

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)"
	var exists bool
	err = db.QueryRow(query, cfg.AdminName).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		hashPassword, err := valuesobject.NewPassword(cfg.AdminPassword)
		if err != nil {
			return nil, err
		}
		if _, err = db.Exec(`INSERT INTO users (id, name, password, role) VALUES ($1, $2, $3, $4)`,
			uuid.New().String(), cfg.AdminName, hashPassword.Value, constants.AdminRole); err != nil {
			return nil, err
		}
	}

	return db, nil
}
