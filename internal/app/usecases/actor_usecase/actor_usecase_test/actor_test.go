package actor_usecase_test

import (
	"context"
	"errors"
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	actorUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/actor_usecase"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestActorUseCase(t *testing.T) {
	actorRepo := mockRepository.NewActorRepository()
	filmRepo := mockRepository.NewFilmRepository()
	useCase := actorUseCase.New(actorRepo, filmRepo)

	testId := uuid.New().String()
	var actorAggr *aggregate.ActorAggregate

	t.Run("Should actor not found", func(t *testing.T) {
		user, _ := useCase.GetById(context.Background(), testId)
		assert.Nil(t, user)
	})

	t.Run("Should create actor", func(t *testing.T) {
		create, err := useCase.Create(context.Background(), appDto.CreateActorUseCaseDto{Name: "Marlen", Gender: "male", Birthday: time.Now().AddDate(-19, 5, 0)})
		assert.Nil(t, err)
		assert.Nil(t, create.Films)
		testId = create.Actor.Id
		user, _ := useCase.GetById(context.Background(), testId)
		assert.NotNil(t, user)
		actorAggr = create
	})

	t.Run("Should incorrect create", func(t *testing.T) {
		create, err := useCase.Create(context.Background(), appDto.CreateActorUseCaseDto{Name: "Marlen", Gender: "incorrect", Birthday: time.Now().AddDate(-19, 5, 0)})
		assert.Nil(t, create)
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusUnprocessableEntity, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}
	})

	t.Run("Should update actor", func(t *testing.T) {
		actorAggr.Actor.Name = "Marlen2"
		update, err := useCase.Update(context.Background(), actorAggr)
		assert.Nil(t, err)
		assert.NotNil(t, update)
		user, _ := useCase.GetById(context.Background(), testId)
		assert.NotNil(t, user)
		assert.Equal(t, "Marlen2", user.Actor.Name)
	})

	t.Run("Should update actor incrrect bad bot found", func(t *testing.T) {
		saveId := actorAggr.Actor.Id
		actorAggr.Actor.Id = "incorrectid"

		update, err := useCase.Update(context.Background(), actorAggr)
		assert.Nil(t, update)
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}

		actorAggr.Actor.Id = saveId
	})

	t.Run("Should delete actor by id", func(t *testing.T) {
		err := useCase.Delete(context.Background(), testId)
		assert.Nil(t, err)
		err = useCase.Delete(context.Background(), testId)
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}
	})

	a, err := useCase.Create(context.Background(), appDto.CreateActorUseCaseDto{Name: "Marlen", Gender: "male", Birthday: time.Now().AddDate(-19, 5, 0)})
	testId = a.Actor.Id

	if err != nil {
		t.Fatal("error create user")
	}

	var (
		film1Id = uuid.New().String()
		film2Id = uuid.New().String()
	)
	_, err = filmRepo.Create(context.Background(), &aggregate.FilmAggregate{Film: model.Film{Id: film1Id, Name: "Titanic", Rate: 10, ReleaseDate: time.Now().AddDate(-13, 0, 0)}})
	if err != nil {
		t.Fatal("error create film")
	}
	_, err = filmRepo.Create(context.Background(), &aggregate.FilmAggregate{Film: model.Film{Id: film2Id, Name: "Hunter", Rate: 9, ReleaseDate: time.Now().AddDate(-10, 0, 0)}})
	if err != nil {
		t.Fatal("error create film")
	}

	t.Run("Should correct add film to actor", func(t *testing.T) {
		err = useCase.AddFilm(context.Background(), testId, film1Id, film2Id)
		assert.Nil(t, err)

		err = useCase.AddFilm(context.Background(), testId)
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusInternalServerError, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}

		err = useCase.AddFilm(context.Background(), testId, "incorrectfilmid")
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}

		err = useCase.AddFilm(context.Background(), "incorrectactorid", film1Id)
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}
	})

	db := inMemDb.New()
	db.CleanUp()
}
