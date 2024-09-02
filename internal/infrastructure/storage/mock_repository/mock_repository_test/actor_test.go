package mock_repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
)

func TestActorRepository(t *testing.T) {
	repo := mockRepository.NewActorRepository()
	db := inMemDb.New()

	birthday, _ := time.Parse("2006-01-02", "1990-01-01")
	actor := &model.Actor{
		Id:       "1",
		Name:     "test_actor",
		Gender:   "male",
		Birthday: birthday,
	}
	actorAggregate := &aggregate.ActorAggregate{Actor: *actor}

	createdActor, err := repo.Create(context.Background(), actorAggregate)
	if err != nil {
		t.Errorf("Ошибка при создании актера: %v", err)
	}
	if createdActor == nil {
		t.Errorf("Созданный актер не должен быть nil")
	}

	retrievedActor, err := repo.GetById(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при получении актера по id: %v", err)
	}
	if retrievedActor == nil || retrievedActor.Actor.Id != "1" {
		t.Errorf("Некорректно полученный актер")
	}

	actorToUpdate := &model.Actor{
		Id:       "1",
		Name:     "updated_actor",
		Gender:   "female",
		Birthday: birthday,
	}
	updatedActor, err := repo.Update(context.Background(), &aggregate.ActorAggregate{Actor: *actorToUpdate})
	if err != nil {
		t.Errorf("Ошибка при обновлении актера: %v", err)
	}
	if updatedActor == nil || updatedActor.Actor.Name != "updated_actor" {
		t.Errorf("Некорректно обновленный актер")
	}

	err = repo.Delete(context.Background(), "1")
	if err != nil {
		t.Errorf("Ошибка при удалении актера: %v", err)
	}

	deletedActor, err := repo.GetById(context.Background(), "1")
	if deletedActor != nil {
		t.Errorf("Актер не был удален")
	}

	film := &model.Film{Id: "1", Name: "titanic", ReleaseDate: time.Now().AddDate(-13, 0, 0), Rate: 10}
	actorToAddFilm := &model.Actor{Id: "1", Name: "actor_with_film"}

	db.Film = append(db.Film, film)
	db.Actor = append(db.Actor, actorToAddFilm)
	err = repo.AddFilm(context.Background(), "1", "1")
	if err != nil {
		t.Errorf("Ошибка при добавлении фильма актеру: %v", err)
	}

	query := domainQuery.ActorRepositoryQuery{
		CurrentPage:    1,
		PageCount:      1,
		WithConnection: []string{},
	}

	actors, _, err := repo.GetByQuery(context.Background(), query)
	if err != nil {
		t.Errorf("Ошибка при выполнении запроса: %v", err)
	}

	if len(actors) != 1 || actors[0].Actor.Id != "1" {
		t.Errorf("Некорректный результат запроса")
	}
}
