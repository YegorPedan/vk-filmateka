package httpv1_test

import (
	"bytes"
	"encoding/json"
	"errors"
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	"github.com/OddEer0/vk-filmoteka/internal/presentation/handlers/httpv1"
	"github.com/OddEer0/vk-filmoteka/internal/presentation/router"
	httpUtils "github.com/OddEer0/vk-filmoteka/pkg/http_utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func initFilmHandler() http.HandlerFunc {
	errHandlerToDefaulHandler := func(next appErrors.AppHandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			err := next(res, req)
			if err != nil {
				var appErr *appErrors.AppError
				if errors.As(err, &appErr) {
					body := appErrors.ResponseError{
						Code:    appErr.Code,
						Message: appErr.Message,
					}
					httpUtils.SendJson(res, appErr.Code, body)
				}
			}
		}
	}
	appHandler := httpv1.NewAppHandlerMock()
	return http.HandlerFunc(errHandlerToDefaulHandler(router.HttpV1RouterFilm(appHandler)))
}

func TestFilmHttpV1Test(t *testing.T) {
	config.MustLoad()
	t.Run("Should create film", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(appDto.CreateFilmUseCaseDto{
			Name:        "Titanic",
			ReleaseDate: time.Now().AddDate(-13, 0, 0),
			Rate:        10,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/film", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should incorrect film", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(appDto.CreateFilmUseCaseDto{
			Name:        "Titanic",
			ReleaseDate: time.Now().AddDate(-13, 0, 0),
			Rate:        11,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/film", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should update film", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(appDto.CreateFilmUseCaseDto{
			Name:        "Titanic",
			ReleaseDate: time.Now().AddDate(-13, 0, 0),
			Rate:        9,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/film", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		var res model.Film
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		if err != nil {
			t.Fatal(err)
		}
		res.Name = "Marvel"
		marshal, err := json.Marshal(res)
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		req, err = http.NewRequest("PUT", "/http/v1/film", bytes.NewBuffer(marshal))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var res2 model.Film
		err = json.Unmarshal(rr.Body.Bytes(), &res2)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Marvel", res2.Name)
	})

	t.Run("Should bad request film", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(appDto.CreateFilmUseCaseDto{
			Name:        "Titanic",
			ReleaseDate: time.Now().AddDate(-13, 0, 0),
			Rate:        9,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/film", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		var res model.Film
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		if err != nil {
			t.Fatal(err)
		}
		res.Rate = 11
		marshal, err := json.Marshal(res)
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		req, err = http.NewRequest("PUT", "/http/v1/film", bytes.NewBuffer(marshal))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should delete film", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(appDto.CreateFilmUseCaseDto{
			Name:        "Titanic",
			ReleaseDate: time.Now().AddDate(-13, 0, 0),
			Rate:        9,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/film", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		var res model.Film
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		if err != nil {
			t.Fatal(err)
		}
		query := "?id=" + res.Id
		rr = httptest.NewRecorder()
		req, err = http.NewRequest("DELETE", "/http/v1/film"+query, nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr = httptest.NewRecorder()
		req, err = http.NewRequest("DELETE", "/http/v1/film", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should error not found film delete", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(appDto.CreateFilmUseCaseDto{
			Name:        "Titanic",
			ReleaseDate: time.Now().AddDate(-13, 0, 0),
			Rate:        9,
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/film", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		var res model.Film
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		if err != nil {
			t.Fatal(err)
		}
		query := "?id=" + "incorrectpassord"
		rr = httptest.NewRecorder()
		req, err = http.NewRequest("DELETE", "/http/v1/film"+query, nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Should get by query", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/http/v1/film", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var res appDto.FilmGetByQueryResult
		err = json.Unmarshal(rr.Body.Bytes(), &res)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 4, len(res.Films))
		assert.Equal(t, 1, res.PageCount)

		rr = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/http/v1/film?page=2&page-count=2&connection=actor&order-by=desc&order-field=rate", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var res2 appDto.FilmGetByQueryResult
		err = json.Unmarshal(rr.Body.Bytes(), &res2)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 2, res2.PageCount)

	})

	t.Run("Should bad request film", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/http/v1/film?page=2&page-count=2&connection=incorrect&order-by=desc&order-field=rate", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/http/v1/film?page=2&page-count=2&connection=actor&order-by=incorrect&order-field=rate", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/http/v1/film?page=2&page-count=2&connection=actor&order-by=desc&order-field=dsa", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/http/v1/film?page=2&page-count=incorrect&connection=actor&order-by=desc&order-field=rate", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/http/v1/film?page=incorrect&page-count=2&connection=actor&order-by=desc&order-field=rate", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should search film", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/http/v1/film/search?search=Marvel", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var res appDto.FilmGetByQueryResult
		err := json.Unmarshal(rr.Body.Bytes(), &res)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(res.Films))
		assert.Equal(t, 1, res.PageCount)
		for _, a := range res.Films {
			t.Log(a.Film.Name)
		}
	})

	t.Run("Should search film bad request", func(t *testing.T) {
		handler := initFilmHandler()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/http/v1/film/search", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/http/v1/film/search?search=Me", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
