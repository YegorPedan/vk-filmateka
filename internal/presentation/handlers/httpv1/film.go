package httpv1

import (
	appDto "github.com/OddEer0/vk-filmoteka/internal/app/app_dto"
	filmUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/film_usecase"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	domainQuery "github.com/OddEer0/vk-filmoteka/internal/domain/repository/domain_query"
	httpUtils "github.com/OddEer0/vk-filmoteka/pkg/http_utils"
	"net/http"
	"strconv"
)

type (
	FilmHandler interface {
		Create(res http.ResponseWriter, req *http.Request) error
		Delete(res http.ResponseWriter, req *http.Request) error
		GetByQuery(res http.ResponseWriter, req *http.Request) error
		Update(res http.ResponseWriter, req *http.Request) error
		SearchByNameAndActorName(res http.ResponseWriter, req *http.Request) error
	}

	filmHandler struct {
		filmUseCase.FilmUseCase
	}
)

func NewFilmHandler(useCase filmUseCase.FilmUseCase) FilmHandler {
	return &filmHandler{
		FilmUseCase: useCase,
	}
}

// @Summary Создание фильма [Админы]
// @Description Доступно только админам
// @Tags film
// @Accept json
// @Produce json
// @Param reg body appDto.CreateFilmUseCaseDto true "Данные фильма"
// @Success 200 {object} model.Film "Данные созданного фильма"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/film [post]
func (f *filmHandler) Create(res http.ResponseWriter, req *http.Request) error {
	var body appDto.CreateFilmUseCaseDto
	err := httpUtils.BodyJson(req, &body)
	if err != nil {
		return appErrors.BadRequest("")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	filmAggregate, err := f.FilmUseCase.Create(req.Context(), body)
	if err != nil {
		return err
	}

	httpUtils.SendJson(res, http.StatusOK, filmAggregate.Film)
	return nil
}

// @Summary Создание фильма [Админы]
// @Description Доступно только админам, ничего не возвоащает
// @Tags film
// @Accept json
// @Produce json
// @Param id query string true "id удаляемого фильма"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/film [delete]
func (f *filmHandler) Delete(res http.ResponseWriter, req *http.Request) error {
	id := req.URL.Query().Get("id")
	if id == "" {
		return appErrors.BadRequest("")
	}

	err := f.FilmUseCase.Delete(req.Context(), id)
	if err != nil {
		return err
	}

	return nil
}

// @Summary Получение фильма
// @Description Можно задавать разные query
// @Tags film
// @Accept json
// @Produce json
// @Param page query string false "текущая страница"
// @Param page-count query string false "кол-во актеров на странице"
// @Param connection query string false "Вернуть вместе со связями (film)"
// @Param order-by query string false "asc либо desc"
// @Param order-field query string false "поле по которому сортируют (rate, name, release_date)"
// @Success 200 {object} appDto.FilmGetByQueryResult "получаемые фильмы"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/film [get]
func (f *filmHandler) GetByQuery(res http.ResponseWriter, req *http.Request) error {
	fQuery := domainQuery.NewFilmRepositoryQuery()
	var err error

	query := req.URL.Query()
	if query.Has("page") {
		fQuery.CurrentPage, err = strconv.Atoi(query.Get("page"))
		if err != nil {
			return appErrors.BadRequest("")
		}
	}
	if query.Has("page-count") {
		fQuery.PageCount, err = strconv.Atoi(query.Get("page-count"))
		if err != nil {
			return appErrors.BadRequest("")
		}
	}
	if query.Has("connection") {
		conn := query.Get("connection")
		if conn != "actor" {
			return appErrors.BadRequest("")
		}
		fQuery.WithConnection = append(fQuery.WithConnection, conn)
	}
	if query.Has("order-by") {
		order := query.Get("order-by")
		switch order {
		case "asc":
			fQuery.OrderBy = domainQuery.Asc
		case "desc":
			fQuery.OrderBy = domainQuery.Desc
		default:
			return appErrors.BadRequest("")
		}
	}
	if query.Has("order-field") {
		field := query.Get("order-field")
		correctFields := []string{"name", "release_date", "rate"}
		has := false
		for _, correctField := range correctFields {
			if field == correctField {
				has = true
			}
		}
		if !has {
			return appErrors.BadRequest("")
		}

		fQuery.SortField = field
	}

	result, err := f.FilmUseCase.GetByQuery(req.Context(), *fQuery)
	if err != nil {
		return err
	}

	httpUtils.SendJson(res, http.StatusOK, result)
	return nil
}

// @Summary Создание фильма [Админы]
// @Description Доступно только админам
// @Tags film
// @Accept json
// @Produce json
// @Param reg body model.Film true "Данные фильма"
// @Success 200 {object} model.Film "Данные обновленного фильма"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/film [put]
func (f *filmHandler) Update(res http.ResponseWriter, req *http.Request) error {
	var body model.Film
	err := httpUtils.BodyJson(req, &body)
	if err != nil {
		return appErrors.BadRequest("")
	}
	defer func() {
		_ = req.Body.Close()
	}()

	filmAggregate, err := aggregate.NewFilmAggregate(body)
	if err != nil {
		return appErrors.UnprocessableEntity("")
	}

	filmAggregate, err = f.FilmUseCase.Update(req.Context(), filmAggregate)
	if err != nil {
		return err
	}

	httpUtils.SendJson(res, http.StatusOK, filmAggregate.Film)
	return nil
}

// @Summary Поиск фильма
// @Description поиск по имени и имени актеров
// @Tags film
// @Accept json
// @Produce json
// @Param search query string false "параметр поиска"
// @Success 200 {object} appDto.FilmGetByQueryResult "получаемые фильмы"
// @Failure 404 {object} appErrors.ResponseError "Ошибка 404"
// @Router /http/v1/film/search [get]
func (f *filmHandler) SearchByNameAndActorName(res http.ResponseWriter, req *http.Request) error {
	query := req.URL.Query()
	if !query.Has("search") {
		return appErrors.BadRequest("search query not found")
	}

	searchedValue := query.Get("search")
	if len(searchedValue) < 3 {
		return appErrors.BadRequest("searched value min 3 chars")
	}

	result, err := f.FilmUseCase.SearchByNameAndActorName(req.Context(), searchedValue)
	if err != nil {
		return err
	}

	httpUtils.SendJson(res, http.StatusOK, result)
	return nil
}
