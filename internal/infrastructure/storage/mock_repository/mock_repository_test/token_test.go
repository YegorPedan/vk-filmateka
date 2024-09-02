package mock_repository_test

import (
	"context"
	"testing"

	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
)

func TestTokenRepository(t *testing.T) {
	repo := mockRepository.NewTokenRepository()
	db := inMemDb.New()

	token := &model.Token{Id: "1", Value: "test_token"}

	createdToken, err := repo.Create(context.Background(), token)
	if err != nil {
		t.Errorf("Ошибка при создании токена: %v", err)
	}
	if createdToken == nil {
		t.Errorf("Созданный токен не должен быть nil")
	}

	retrievedToken, err := repo.GetById(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при получении токена по id: %v", err)
	}
	if retrievedToken == nil || retrievedToken.Id != "1" {
		t.Errorf("Некорректно полученный токен")
	}

	tokenToUpdate := &model.Token{Id: "1", Value: "updated_token"}
	updatedToken, err := repo.Update(context.Background(), tokenToUpdate)
	if err != nil {
		t.Errorf("Ошибка при обновлении токена: %v", err)
	}
	if updatedToken == nil || updatedToken.Value != "updated_token" {
		t.Errorf("Некорректно обновленный токен")
	}

	err = repo.Delete(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при удалении токена: %v", err)
	}

	deletedToken, err := repo.GetById(context.Background(), "1")
	if deletedToken != nil {
		t.Errorf("Токен не был удален")
	}

	tokenToDeleteByValue := &model.Token{Id: "2", Value: "delete_by_value_token"}
	db.Tokens = append(db.Tokens, tokenToDeleteByValue)

	err = repo.DeleteByValue(context.Background(), "delete_by_value_token")
	if err != nil {
		t.Errorf("Ошибка при удалении токена по значению: %v", err)
	}

	hasTokenByValue, err := repo.HasByValue(context.Background(), "delete_by_value_token")
	if hasTokenByValue {
		t.Errorf("Токен по значению не был удален")
	}
}
