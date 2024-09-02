package httpv1_test

import (
	"bytes"
	"encoding/json"
	"errors"
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
)

func initActorHandler() http.HandlerFunc {
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
	return http.HandlerFunc(errHandlerToDefaulHandler(router.HttpV1RouterActor(appHandler)))
}

func TestActorHttpV1Test(t *testing.T) {
	cfg := config.MustLoad()
	t.Run("Should create actor", func(t *testing.T) {
		handler := initActorHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Jason",
			"gender":   "male",
			"birthday": "2004-03-17T18:43:48.52645+03:00",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/actor", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should create actor bad request", func(t *testing.T) {
		handler := initActorHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Jason",
			"gender":   "dsa",
			"birthday": "2004-03-17T18:43:48.52645+03:00",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/actor", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should update actor", func(t *testing.T) {
		handler := initActorHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Scala",
			"gender":   "male",
			"birthday": "2004-03-17T18:43:48.52645+03:00",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/actor", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var body model.Actor
		_ = json.Unmarshal(rr.Body.Bytes(), &body)

		body.Name = "Python"
		re, _ := json.Marshal(body)
		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", "/http/v1/actor", bytes.NewBuffer(re))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		body.Id = "dsadsa"
		re, _ = json.Marshal(body)
		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", "/http/v1/actor", bytes.NewBuffer(re))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		requestBody, _ = json.Marshal(map[string]string{
			"name":     "Scala",
			"gender":   "dsadsad",
			"birthday": "2004-03-17T18:43:48.52645+03:00",
		})
		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", "/http/v1/actor", bytes.NewBuffer(requestBody))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should delete actor", func(t *testing.T) {
		handler := initActorHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Scala",
			"gender":   "male",
			"birthday": "2004-03-17T18:43:48.52645+03:00",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/actor", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var body model.Film
		_ = json.Unmarshal(rr.Body.Bytes(), &body)
		id := body.Id

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/http/v1/actor"+"?id="+id, bytes.NewBuffer(requestBody))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/http/v1/actor"+"?id="+id, bytes.NewBuffer(requestBody))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("DELETE", "/http/v1/actor", bytes.NewBuffer(requestBody))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should get by query", func(t *testing.T) {
		handler := initActorHandler()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/http/v1/actor?page=1&page-count=2&connection=film", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should bad request get by query", func(t *testing.T) {
		handler := initActorHandler()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/http/v1/actor?page=1&page-count=2&connection=fil", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		rr = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/http/v1/actor?page=1&page-count=dsa&connection=film", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		req, _ = http.NewRequest("GET", "/http/v1/actor?page=adsads&page-count=2&connection=film", nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	cfg.Env = "dev"
	t.Run("Should create actor", func(t *testing.T) {
		handler := initActorHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Jason",
			"gender":   "male",
			"birthday": "2004-03-17T18:43:48.52645+03:00",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/actor", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should create actor", func(t *testing.T) {
		authHandler := initAppHandler()
		rrAuth := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Admin",
			"password": "Adminadmin41",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/login", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		authHandler.ServeHTTP(rrAuth, req)

		handler := initActorHandler()
		rr := httptest.NewRecorder()
		requestBody, err = json.Marshal(map[string]string{
			"name":     "Jason",
			"gender":   "male",
			"birthday": "2004-03-17T18:43:48.52645+03:00",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err = http.NewRequest("POST", "/http/v1/actor", bytes.NewBuffer(requestBody))
		setToken(rrAuth, req)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
	cfg.Env = "test"
}
