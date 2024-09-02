package postgresRepository

import (
	"context"
	"database/sql"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
)

type userRepository struct {
	db *sql.DB
}

func (u userRepository) Create(ctx context.Context, userAggregate *aggregate.UserAggregate) (*aggregate.UserAggregate, error) {
	query := "INSERT INTO users (id, name, password, role) VALUES ($1, $2, $3, $4) RETURNING id, name, password, role"
	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	var user model.User
	err = stmt.QueryRowContext(ctx, userAggregate.User.Id, userAggregate.User.Name, userAggregate.User.Password.Value, userAggregate.User.Role).Scan(&user.Id, &user.Name, &user.Password.Value, &user.Role)
	if err != nil {
		return nil, err
	}

	userAggregate.User = user
	return userAggregate, nil
}

func (u userRepository) Update(ctx context.Context, userAggregate *aggregate.UserAggregate) (*aggregate.UserAggregate, error) {
	query := "UPDATE users SET name = $1, password = $2, role = $3 WHERE id = $4 RETURNING id, name, password, role"
	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	var user model.User
	err = stmt.QueryRowContext(ctx, userAggregate.User.Name, userAggregate.User.Password.Value, userAggregate.User.Role, userAggregate.User.Id).Scan(&user.Id, &user.Name, &user.Password.Value, &user.Role)
	if err != nil {
		return nil, err
	}

	userAggregate.User = user
	return userAggregate, nil
}

func (u userRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := u.db.ExecContext(ctx, query, id)
	return err
}

func (u userRepository) GetById(ctx context.Context, id string) (*aggregate.UserAggregate, error) {
	query := "SELECT id, name, password, role FROM users WHERE id = $1"
	row := u.db.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Password.Value, &user.Role)
	if err != nil {
		return nil, err
	}

	return &aggregate.UserAggregate{User: user}, nil
}

func (u userRepository) HasUserByName(ctx context.Context, name string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)"
	var exists bool
	err := u.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u userRepository) GetByName(ctx context.Context, name string) (*aggregate.UserAggregate, error) {
	query := "SELECT id, name, password, role FROM users WHERE name = $1"
	row := u.db.QueryRowContext(ctx, query, name)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Password.Value, &user.Role)
	if err != nil {
		return nil, err
	}

	return &aggregate.UserAggregate{User: user}, nil
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}
