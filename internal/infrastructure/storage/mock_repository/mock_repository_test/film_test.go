package mock_repository_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
)

func TestFilmRepository(t *testing.T) {
	// Инициализируем репозиторий и базу данных
	repo := mockRepository.NewFilmRepository()
	db := inMemDb.New()

	// Создаем тестовые данные
	desc := "Test film description"
	releaseDate, _ := time.Parse("2006-01-02", "2022-01-01")
	film := &model.Film{
		Id:          "100",
		Name:        "test__film",
		ReleaseDate: releaseDate,
		Rate:        5,
		Description: &desc,
	}
	filmAggregate := &aggregate.FilmAggregate{Film: *film}

	// Тест метода Create
	createdFilm, err := repo.Create(context.Background(), filmAggregate)
	if err != nil {
		t.Errorf("Ошибка при создании фильма: %v", err)
	}
	if createdFilm == nil {
		t.Errorf("Созданный фильм не должен быть nil")
	}

	// Тест метода GetById
	retrievedFilm, err := repo.GetById(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при получении фильма по ID: %v", err)
	}
	if retrievedFilm == nil || retrievedFilm.Film.Id != "1" {
		t.Errorf("Некорректно полученный фильм")
	}

	desc2 := "Updated film description"
	filmToUpdate := &model.Film{
		Id:          "1",
		Name:        "updated_film",
		ReleaseDate: releaseDate,
		Rate:        5,
		Description: &desc2,
	}
	updatedFilm, err := repo.Update(context.Background(), &aggregate.FilmAggregate{Film: *filmToUpdate})
	if err != nil {
		t.Errorf("Ошибка при обновлении фильма: %v", err)
	}
	if updatedFilm == nil || updatedFilm.Film.Name != "updated_film" {
		t.Errorf("Некорректно обновленный фильм")
	}

	// Тест метода Delete
	err = repo.Delete(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при удалении фильма: %v", err)
	}

	// Проверяем, что фильм действительно удален
	deletedFilm, err := repo.GetById(context.Background(), "1")
	if deletedFilm != nil {
		t.Errorf("Фильм не был удален")
	}

	// Тест метода SearchByNameAndActorName
	// Создаем несколько тестовых фильмов
	film1 := &model.Film{Id: "12121", Name: "superpuper1"}
	film2 := &model.Film{Id: "2323", Name: "superpuper3"}
	film3 := &model.Film{Id: "34242", Name: "superpuper2"}

	db.Film = append(db.Film, film1, film2, film3)

	searchValue := "superpuper"
	foundFilms, _, err := repo.SearchByNameAndActorName(context.Background(), searchValue)
	if err != nil {
		t.Errorf("Ошибка при поиске фильмов: %v", err)
	}
	if len(foundFilms) != 3 {
		t.Errorf("Некорректное количество найденных фильмов")
	}

	// Проверяем, что все найденные фильмы содержат искомое значение в имени
	for _, f := range foundFilms {
		if !strings.Contains(f.Film.Name, searchValue) {
			t.Errorf("Найденный фильм не содержит искомое значение в имени")
		}
	}
}
