package mockRepository

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/repository"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
)

type userRepository struct {
	db *inMemDb.InMemDb
}

func (u userRepository) Create(ctx context.Context, data *aggregate.UserAggregate) (*aggregate.UserAggregate, error) {
	has := slices.ContainsFunc(u.db.Users, func(item *model.User) bool {
		if item.Name == data.User.Name {
			return true
		}
		return false
	})

	if has {
		return nil, errors.New("conflict fields")
	}

	u.db.Users = append(u.db.Users, &data.User)

	if data.Token != nil {
		token := model.Token{Id: data.Token.Id, Value: data.Token.Value}
		u.db.Tokens = append(u.db.Tokens, &token)
	}

	return data, nil
}

func (u userRepository) Update(ctx context.Context, data *aggregate.UserAggregate) (*aggregate.UserAggregate, error) {
	has := false

	for i, user := range u.db.Users {
		if data.User.Id == user.Id {
			has = true
			copyUser := data.User
			u.db.Users[i] = &copyUser
			break
		}
	}

	if has {
		return data, nil
	}
	return nil, sql.ErrNoRows
}

func (u userRepository) Delete(ctx context.Context, id string) error {
	has := false

	var filteredUsers []*model.User

	for _, user := range u.db.Users {
		if user.Id != id {
			filteredUsers = append(filteredUsers, user)
		} else {
			has = true
		}
	}

	u.db.Users = filteredUsers

	if has {
		return nil
	}
	return sql.ErrNoRows
}

func (u userRepository) GetById(ctx context.Context, id string) (*aggregate.UserAggregate, error) {
	var searched *model.User = nil
	for _, user := range u.db.Users {
		if user.Id == id {
			searched = user
		}
	}
	if searched != nil {
		return &aggregate.UserAggregate{User: *searched}, nil
	}
	return nil, sql.ErrNoRows
}

func (u userRepository) HasUserByName(ctx context.Context, name string) (bool, error) {
	return slices.ContainsFunc(u.db.Users, func(item *model.User) bool {
		if item.Name == name {
			return true
		}
		return false
	}), nil
}

func (u userRepository) GetByName(ctx context.Context, name string) (*aggregate.UserAggregate, error) {
	for _, user := range u.db.Users {
		if user.Name == name {
			return &aggregate.UserAggregate{User: *user}, nil
		}
	}
	return nil, sql.ErrNoRows
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{inMemDb.New()}
}
