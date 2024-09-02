package repository

import (
	"context"

	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
)

type TokenRepository interface {
	Create(ctx context.Context, model *model.Token) (*model.Token, error)
	Update(ctx context.Context, model *model.Token) (*model.Token, error)
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (*model.Token, error)
	DeleteByValue(ctx context.Context, value string) error
	HasByValue(ctx context.Context, value string) (bool, error)
}
