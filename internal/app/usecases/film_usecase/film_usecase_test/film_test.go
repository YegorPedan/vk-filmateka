package film_usecase_test

import (
	"context"
	"errors"
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	filmUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/film_usecase"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	inMemDb "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/in_mem_db"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestFilmUseCase(t *testing.T) {
	filmRepo := mockRepository.NewFilmRepository()
	useCase := filmUseCase.New(filmRepo)

	testId := ""
	var film *aggregate.FilmAggregate
	t.Run("Should create film", func(t *testing.T) {
		create, err := useCase.Create(context.Background(), appDto.CreateFilmUseCaseDto{Name: "Titanic", ReleaseDate: time.Now().AddDate(-13, 0, 0), Rate: 10})
		assert.Nil(t, err)
		assert.NotNil(t, create)

		testId = create.Film.Id
		id, err := useCase.GetById(context.Background(), testId)
		assert.Nil(t, err)
		assert.Equal(t, testId, id.Film.Id)
		film = create
	})

	t.Run("Should update film", func(t *testing.T) {
		film.Film.Name = "Titanic 2"
		update, err := useCase.Update(context.Background(), film)
		assert.Nil(t, err)
		assert.Equal(t, "Titanic 2", update.Film.Name)

		film.Film.Id = "incorrect"
		update, err = useCase.Update(context.Background(), film)
		assert.Nil(t, update)
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}

		film.Film.Id = testId
	})

	t.Run("Should delete film", func(t *testing.T) {
		err := useCase.Delete(context.Background(), testId)
		assert.Nil(t, err)
		film, err := useCase.GetById(context.Background(), testId)
		assert.Nil(t, film)
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}

		err = useCase.Delete(context.Background(), testId)
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}
	})

	film, err := useCase.Create(context.Background(), appDto.CreateFilmUseCaseDto{Name: "Titanic", ReleaseDate: time.Now().AddDate(-13, 0, 0), Rate: 10})
	if err != nil {
		t.Fatal("incorrect create")
	}

	t.Run("Should search by name and actor name", func(t *testing.T) {
		res, err := useCase.SearchByNameAndActorName(context.Background(), "Tita")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.Films))
		assert.Equal(t, 1, res.PageCount)
		res, err = useCase.SearchByNameAndActorName(context.Background(), "Murmur")
		assert.Nil(t, res)
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) {
			assert.Equal(t, http.StatusNotFound, appErr.Code)
		} else {
			t.Fatal("incorrect error type")
		}
	})

	db := inMemDb.New()
	db.CleanUp()
}
