package httpv1

import (
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	actorUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/actor_usecase"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	"github.com/OddEer0/vk-filmoteka/internal/presentation/dto"
	httpUtils "github.com/OddEer0/vk-filmoteka/pkg/http_utils"
	"net/http"
	"strconv"
)

type (
	ActorHandler interface {
		Create(res http.ResponseWriter, req *http.Request) error
		Delete(res http.ResponseWriter, req *http.Request) error
		GetByQuery(res http.ResponseWriter, req *http.Request) error
		Update(res http.ResponseWriter, req *http.Request) error
		AddFilm(res http.ResponseWriter, req *http.Request) error
	}

	actorHandler struct {
		actorUseCase.ActorUseCase
	}
)

func NewActorHandler(useCase actorUseCase.ActorUseCase) ActorHandler {
	return &actorHandler{
		ActorUseCase: useCase,
	}
}

// @Summary Создание актера [Админы]
// @Description Доступно только админам
// @Tags actor
// @Accept json
// @Produce json
// @Param reg body appDto.CreateActorUseCaseDto true "Данные актера"
// @Success 200 {object} model.Actor "Данные созданного актера"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/actor [post]
func (a *actorHandler) Create(res http.ResponseWriter, req *http.Request) error {
	var body appDto.CreateActorUseCaseDto
	err := httpUtils.BodyJson(req, &body)
	if err != nil {
		return appErrors.BadRequest("")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	actorAggregate, err := a.ActorUseCase.Create(req.Context(), body)
	if err != nil {
		return err
	}

	httpUtils.SendJson(res, http.StatusOK, actorAggregate.Actor)
	return nil
}

// @Summary Удаление актера [Админы]
// @Description Доступно только админам, ничего ответом не возвращает
// @Tags actor
// @Accept json
// @Produce json
// @Param id query string true "id удаляемого пользователья"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/actor [delete]
func (a *actorHandler) Delete(res http.ResponseWriter, req *http.Request) error {
	id := req.URL.Query().Get("id")
	if id == "" {
		return appErrors.BadRequest("")
	}

	err := a.ActorUseCase.Delete(req.Context(), id)
	if err != nil {
		return err
	}
	return nil
}

// @Summary Получение актера
// @Description Можно задавать разные query
// @Tags actor
// @Accept json
// @Produce json
// @Param page query string false "текущая страница"
// @Param page-count query string false "кол-во актеров на странице"
// @Param connection query string false "Вернуть вместе со связями (film)"
// @Success 200 {object} appDto.ActorGetByQueryResult "получаемые актеры"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/actor [get]
func (a *actorHandler) GetByQuery(res http.ResponseWriter, req *http.Request) error {
	query := req.URL.Query()
	fQuery := domainQuery.NewActorRepositoryQuery()
	var err error

	if query.Has("page") {
		currentPageQ := req.URL.Query().Get("page")
		fQuery.CurrentPage, err = strconv.Atoi(currentPageQ)
		if err != nil {
			return appErrors.BadRequest("invalid page")
		}
	}
	if query.Has("page-count") {
		pageCountQ := req.URL.Query().Get("page-count")
		fQuery.PageCount, err = strconv.Atoi(pageCountQ)
		if err != nil {
			return appErrors.BadRequest("invalid page count")
		}
	}
	if query.Has("connection") {
		connectionQ := req.URL.Query().Get("connection")
		if connectionQ != "film" {
			return appErrors.BadRequest("invalid connection")
		}
		fQuery.WithConnection = append(fQuery.WithConnection, connectionQ)
	}

	result, err := a.ActorUseCase.GetByQuery(req.Context(), *fQuery)
	if err != nil {
		return err
	}

	httpUtils.SendJson(res, http.StatusOK, result)

	return nil
}

// @Summary Обновление актера [Админы]
// @Description Доступно только админам
// @Tags actor
// @Accept json
// @Produce json
// @Param reg body model.Actor true "Данные актера"
// @Success 200 {object} model.Actor "Данные актера"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/actor [put]
func (a *actorHandler) Update(res http.ResponseWriter, req *http.Request) error {
	var body model.Actor
	err := httpUtils.BodyJson(req, &body)
	if err != nil {
		return appErrors.BadRequest("")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	actorAggregate, err := aggregate.NewActorAggregate(body)
	if err != nil {
		return appErrors.UnprocessableEntity("")
	}

	actorAggregate, err = a.ActorUseCase.Update(req.Context(), actorAggregate)
	if err != nil {
		return err
	}

	httpUtils.SendJson(res, http.StatusOK, actorAggregate.Actor)
	return nil
}

// @Summary Обновление актера [Админы]
// @Description Доступно только админам
// @Tags actor
// @Accept json
// @Produce json
// @Param reg body dto.AddFilmToActorDto true "Данные id для связывания актера и фильма"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/actor/add-film [post]
func (a *actorHandler) AddFilm(res http.ResponseWriter, req *http.Request) error {
	var body dto.AddFilmToActorDto
	err := httpUtils.BodyJson(req, &body)
	if err != nil {
		return appErrors.BadRequest("", "Target: Handler")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	err = a.ActorUseCase.AddFilm(req.Context(), body.ActorId, body.FilmIds...)
	if err != nil {
		return err
	}

	return nil
}
