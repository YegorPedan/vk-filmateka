package postgresRepository

import (
	"context"
	"database/sql"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
)

type tokenRepository struct {
	db *sql.DB
}

func (t tokenRepository) Create(ctx context.Context, token *model.Token) (*model.Token, error) {
	query := "INSERT INTO tokens (id, value) VALUES ($1, $2) RETURNING id, value"
	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	err = stmt.QueryRowContext(ctx, token.Id, token.Value).Scan(&token.Id, &token.Value)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t tokenRepository) Update(ctx context.Context, token *model.Token) (*model.Token, error) {
	query := "UPDATE tokens SET value = $1 WHERE id = $2 RETURNING id, value"
	stmt, err := t.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	err = stmt.QueryRowContext(ctx, token.Value, token.Id).Scan(&token.Id, &token.Value)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t tokenRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM tokens WHERE id = $1"
	_, err := t.db.ExecContext(ctx, query, id)
	return err
}

func (t tokenRepository) GetById(ctx context.Context, id string) (*model.Token, error) {
	query := "SELECT id, value FROM tokens WHERE id = $1"
	row := t.db.QueryRowContext(ctx, query, id)

	var token model.Token
	err := row.Scan(&token.Id, &token.Value)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t tokenRepository) DeleteByValue(ctx context.Context, value string) error {
	query := "DELETE FROM tokens WHERE value = $1"
	_, err := t.db.ExecContext(ctx, query, value)
	return err
}

func (t tokenRepository) HasByValue(ctx context.Context, value string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tokens WHERE value = $1)"
	var exists bool
	err := t.db.QueryRowContext(ctx, query, value).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func NewTokenRepository(db *sql.DB) repository.TokenRepository {
	return &tokenRepository{db: db}
}
