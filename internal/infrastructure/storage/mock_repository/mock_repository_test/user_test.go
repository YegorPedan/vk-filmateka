package mock_repository_test

import (
	"context"
	"testing"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
)

func TestUserRepository(t *testing.T) {
	repo := mockRepository.NewUserRepository()

	user := &model.User{Id: "1", Name: "test_user"}
	userAggregate := &aggregate.UserAggregate{User: *user}

	createdUser, err := repo.Create(context.Background(), userAggregate)
	if err != nil {
		t.Errorf("Ошибка при создании пользователя: %v", err)
	}
	if createdUser == nil {
		t.Errorf("Созданный пользователь не должен быть nil")
	}

	retrievedUser, err := repo.GetById(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при получении пользователя по id: %v", err)
	}
	if retrievedUser == nil || retrievedUser.User.Id != "1" {
		t.Errorf("Некорректно полученный пользователь")
	}

	hasUser, err := repo.HasUserByName(context.Background(), "test_user")
	if err != nil {
		t.Errorf("Ошибка при проверке существования пользователя по имени: %v", err)
	}
	if !hasUser {
		t.Errorf("Ожидался пользователь с именем 'test_user', но не был найден")
	}

	retrievedByName, err := repo.GetByName(context.Background(), "test_user")
	if err != nil {
		t.Errorf("Ошибка при получении пользователя по имени: %v", err)
	}
	if retrievedByName == nil || retrievedByName.User.Name != "test_user" {
		t.Errorf("Некорректно полученный пользователь по имени")
	}

	userToUpdate := &model.User{Id: "1", Name: "updated_user"}
	updatedUser, err := repo.Update(context.Background(), &aggregate.UserAggregate{User: *userToUpdate})
	if err != nil {
		t.Errorf("Ошибка при обновлении пользователя: %v", err)
	}
	if updatedUser == nil || updatedUser.User.Name != "updated_user" {
		t.Errorf("Некорректно обновленный пользователь")
	}

	err = repo.Delete(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при удалении пользователя: %v", err)
	}

	deletedUser, err := repo.GetById(context.Background(), "1")
	if deletedUser != nil {
		t.Errorf("Пользователь не был удален")
	}
}
