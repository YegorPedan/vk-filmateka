package httpv1_test

import (
	"bytes"
	"encoding/json"
	"errors"
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	"github.com/OddEer0/vk-filmoteka/internal/presentation/handlers/httpv1"
	"github.com/OddEer0/vk-filmoteka/internal/presentation/router"
	httpUtils "github.com/OddEer0/vk-filmoteka/pkg/http_utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initAppHandler() http.HandlerFunc {
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
	return http.HandlerFunc(errHandlerToDefaulHandler(router.HttpV1RouterAuth(appHandler)))
}

func setToken(rr *httptest.ResponseRecorder, req *http.Request) {
	cookies := rr.Result().Cookies()
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
}

func TestAuthHandler(t *testing.T) {
	config.MustLoad()
	t.Run("Should registration", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Marlens",
			"password": "c21312121314",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/registration", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var body appDto.ResponseUserDto
		err = json.Unmarshal(rr.Body.Bytes(), &body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Marlens", body.Name)
	})

	t.Run("Should incorrect registration", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()

		requestBody, err := json.Marshal(map[string]string{
			"name":     "Marlens",
			"password": "c21312121314",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/registration", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusConflict, rr.Code)

	})

	t.Run("Should bad request password registration", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()

		requestBody, err := json.Marshal(map[string]string{
			"name":     "Marlens",
			"password": "incorrect",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/registration", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

	})

	t.Run("Should bad request name registration", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()

		requestBody, err := json.Marshal(map[string]string{
			"name":     "Ma",
			"password": "Incorrect32",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/registration", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

	})

	t.Run("Should login", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Marlens",
			"password": "c21312121314",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/login", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var body appDto.ResponseUserDto
		err = json.Unmarshal(rr.Body.Bytes(), &body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Marlens", body.Name)
	})

	t.Run("Should incorrect login", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Marlenss",
			"password": "c21312121314",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/login", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("Should logout", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Marlens",
			"password": "c21312121314",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/login", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		rr2 := httptest.NewRecorder()
		req2, err := http.NewRequest("GET", "/http/v1/auth/logout", nil)
		if err != nil {
			t.Fatal(err)
		}
		setToken(rr, req2)
		handler.ServeHTTP(rr2, req2)
		assert.Equal(t, http.StatusOK, rr2.Code)

		rr2 = httptest.NewRecorder()
		req2, err = http.NewRequest("GET", "/http/v1/auth/logout", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr2, req2)
		assert.Equal(t, http.StatusBadRequest, rr2.Code)
	})

	t.Run("Should refresh", func(t *testing.T) {
		handler := initAppHandler()
		rr := httptest.NewRecorder()
		requestBody, err := json.Marshal(map[string]string{
			"name":     "Marlens",
			"password": "c21312121314",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/http/v1/auth/login", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		rr2 := httptest.NewRecorder()
		req2, err := http.NewRequest("POST", "/http/v1/auth/refresh", nil)
		if err != nil {
			t.Fatal(err)
		}
		setToken(rr, req2)
		handler.ServeHTTP(rr2, req2)
		assert.Equal(t, http.StatusOK, rr2.Code)
	})
}
