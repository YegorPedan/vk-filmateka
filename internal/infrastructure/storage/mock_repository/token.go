package mockRepository

import (
	"context"
	"database/sql"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
)

type tokenRepository struct {
	db *inMemDb.InMemDb
}

func (t tokenRepository) Create(ctx context.Context, data *model.Token) (*model.Token, error) {
	token := model.Token{Id: data.Id, Value: data.Value}
	t.db.Tokens = append(t.db.Tokens, &token)
	return &token, nil
}

func (t tokenRepository) Update(ctx context.Context, data *model.Token) (*model.Token, error) {
	for _, item := range t.db.Tokens {
		if item.Id == data.Id {
			item.Value = data.Value
			return item, nil
		}
	}
	return nil, nil
}

func (t tokenRepository) Delete(ctx context.Context, id string) error {
	var filteredTokens []*model.Token

	for _, item := range t.db.Tokens {
		if item.Id != id {
			filteredTokens = append(filteredTokens, item)
		}
	}

	t.db.Tokens = filteredTokens

	return nil
}

func (t tokenRepository) GetById(ctx context.Context, id string) (*model.Token, error) {
	var token *model.Token = nil
	found := false

	for _, item := range t.db.Tokens {
		if item.Id == id {
			token = item
			found = true
			break
		}
	}

	if !found {
		return nil, sql.ErrNoRows
	}
	return token, nil
}

func (t tokenRepository) DeleteByValue(ctx context.Context, value string) error {
	var filteredTokens []*model.Token

	for _, item := range t.db.Tokens {
		if item.Value != value {
			filteredTokens = append(filteredTokens, item)
		}
	}

	t.db.Tokens = filteredTokens

	return nil
}

func (t tokenRepository) HasByValue(ctx context.Context, value string) (bool, error) {
	for _, token := range t.db.Tokens {
		if token.Value == value {
			return true, nil
		}
	}
	return false, nil
}

func NewTokenRepository() repository.TokenRepository {
	return &tokenRepository{inMemDb.New()}
}
