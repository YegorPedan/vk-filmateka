package httpv1

import (
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	authUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/auth_usecase"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	httpUtils "github.com/OddEer0/vk-filmoteka/pkg/http_utils"
	"net/http"
	"time"

	_ "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	_ "github.com/OddEer0/vk-filmoteka/internal/presentation/dto"
)

type (
	AuthHandler interface {
		Registration(res http.ResponseWriter, req *http.Request) error
		Login(res http.ResponseWriter, req *http.Request) error
		Logout(res http.ResponseWriter, req *http.Request) error
		Refresh(res http.ResponseWriter, req *http.Request) error
	}

	authHandler struct {
		authUseCase.AuthUseCase
	}
)

func NewAuthHandler(useCase authUseCase.AuthUseCase) AuthHandler {
	return &authHandler{
		AuthUseCase: useCase,
	}
}

// @Summary Регистрация пользователя
// @Description Ответом при успешном регистрация получаем свои данные
// @Tags auth
// @Accept json
// @Produce json
// @Param reg body appDto.RegistrationUseCaseDto true "Данные нового пользователя"
// @Success 200 {object} appDto.ResponseUserDto "Данные созданного пользователя"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/auth/registration [post]
func (a *authHandler) Registration(res http.ResponseWriter, req *http.Request) error {
	var body appDto.RegistrationUseCaseDto
	err := httpUtils.BodyJson(req, &body)
	if err != nil {
		return appErrors.BadRequest("")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	registerResult, err := a.AuthUseCase.Registration(req.Context(), body)
	if err != nil {
		return err
	}

	err = a.setToken(res, registerResult.Tokens.RefreshToken, registerResult.Tokens.AccessToken)
	if err != nil {
		return appErrors.InternalServerError(err.Error())
	}
	httpUtils.SendJson(res, http.StatusOK, registerResult.User)
	return nil
}

// @Summary Логин пользователя
// @Description Ответом при успешном Логине получаем свои данные
// @Tags auth
// @Accept json
// @Produce json
// @Param login body appDto.LoginUseCaseDto true "Данные пользователя"
// @Success 200 {object} appDto.ResponseUserDto "Данные пользователя"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/auth/login [post]
func (a *authHandler) Login(res http.ResponseWriter, req *http.Request) error {
	var body appDto.LoginUseCaseDto
	err := httpUtils.BodyJson(req, &body)
	if err != nil {
		return appErrors.BadRequest("")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	loginResult, err := a.AuthUseCase.Login(req.Context(), body)
	if err != nil {
		return err
	}

	err = a.setToken(res, loginResult.Tokens.RefreshToken, loginResult.Tokens.AccessToken)
	if err != nil {
		return appErrors.InternalServerError("set token error")
	}

	httpUtils.SendJson(res, http.StatusOK, loginResult.User)
	return nil
}

// @Summary Обновление access токена пользователя
// @Description Ответом ничего не получает. Чистится куки
// @Tags auth
// @Accept json
// @Produce json
// @CookieParam accessToken string true "Идентификатор сессии"
// @CookieParam refreshToken string true "Идентификатор сессии"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/auth/logout [get]
func (a *authHandler) Logout(res http.ResponseWriter, req *http.Request) error {
	token, err := req.Cookie("refreshToken")
	if err != nil {
		return appErrors.BadRequest("")
	}
	err = a.AuthUseCase.Logout(req.Context(), token.Value)
	if err != nil {
		return err
	}

	a.removeToken(res)
	return nil
}

// @Summary Обновление access токена пользователя
// @Description Ответом при успешном Логине получаем свои данные
// @Tags auth
// @Accept json
// @Produce json
// @CookieParam refreshToken string true "Идентификатор сессии"
// @Success 200 {object} appDto.ResponseUserDto "Данные созданного пользователя"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/auth/refresh [post]
func (a *authHandler) Refresh(res http.ResponseWriter, req *http.Request) error {
	token, err := req.Cookie("refreshToken")
	if err != nil {
		return appErrors.BadRequest("")
	}
	result, err := a.AuthUseCase.Refresh(req.Context(), token.Value)
	if err != nil {
		return err
	}
	err = a.setToken(res, result.Tokens.RefreshToken, result.Tokens.AccessToken)
	if err != nil {
		return appErrors.InternalServerError("set token error")
	}
	httpUtils.SendJson(res, http.StatusOK, result.User)
	return nil
}

func (a *authHandler) setToken(res http.ResponseWriter, refreshToken string, accessToken string) error {
	cfg := config.NewConfig()
	refreshTokenTime, err := time.ParseDuration(cfg.RefreshTokenTime)
	if err != nil {
		return err
	}
	accessTokenTime, err := time.ParseDuration(cfg.AccessTokenTime)
	if err != nil {
		return err
	}

	refreshCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(refreshTokenTime.Minutes()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	accessCookie := http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(accessTokenTime.Minutes()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(res, &accessCookie)
	http.SetCookie(res, &refreshCookie)
	return nil
}

func (a *authHandler) removeToken(res http.ResponseWriter) {
	refreshCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	accessCookie := http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(res, &accessCookie)
	http.SetCookie(res, &refreshCookie)
}
