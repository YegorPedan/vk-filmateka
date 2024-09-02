package router

import (
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
	"github.com/OddEer0/vk-filmoteka/internal/presentation/middleware"
	"net/http"
	"strings"

	"github.com/OddEer0/vk-filmoteka/internal/presentation/handlers/httpv1"
)

const (
	HttpV1Prefix = "/http/v1"
)

func HttpV1Router(appHandler *httpv1.AppHandler) appErrors.AppHandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) error {
		path := strings.TrimPrefix(req.URL.Path, HttpV1Prefix)

		switch {
		case strings.HasPrefix(path, "/auth"):
			return HttpV1RouterAuth(appHandler)(res, req)
		case strings.HasPrefix(path, "/actor"):
			return HttpV1RouterActor(appHandler)(res, req)
		case strings.HasPrefix(path, "/film"):
			return HttpV1RouterFilm(appHandler)(res, req)
		default:
			http.NotFound(res, req)
		}
		return nil
	}
}

func HttpV1RouterAuth(appHandler *httpv1.AppHandler) appErrors.AppHandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) error {
		path := strings.TrimPrefix(req.URL.Path, HttpV1Prefix+"/auth")
		switch {
		case req.Method == http.MethodPost && path == "/registration":
			return appHandler.AuthHandler.Registration(res, req)
		case req.Method == http.MethodPost && path == "/login":
			return appHandler.AuthHandler.Login(res, req)
		case req.Method == http.MethodPost && path == "/refresh":
			return appHandler.AuthHandler.Refresh(res, req)
		case req.Method == http.MethodGet && path == "/logout":
			return appHandler.AuthHandler.Logout(res, req)
		default:
			http.NotFound(res, req)
		}
		return nil
	}
}

func HttpV1RouterActor(appHandler *httpv1.AppHandler) appErrors.AppHandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) error {
		path := strings.TrimPrefix(req.URL.Path, HttpV1Prefix+"/actor")

		adminMiddleware := middleware.AuthRoleMiddleware(constants.AdminRole)
		switch {
		case http.MethodPost == req.Method && path == "/add-film":
			return adminMiddleware(appHandler.ActorHandler.AddFilm)(res, req)
		case http.MethodGet == req.Method:
			return appHandler.ActorHandler.GetByQuery(res, req)
		case http.MethodPost == req.Method:
			return adminMiddleware(appHandler.ActorHandler.Create)(res, req)
		case http.MethodPut == req.Method:
			return adminMiddleware(appHandler.ActorHandler.Update)(res, req)
		case http.MethodDelete == req.Method:
			return adminMiddleware(appHandler.ActorHandler.Delete)(res, req)
		default:
			http.NotFound(res, req)
		}
		return nil
	}
}

func HttpV1RouterFilm(appHandler *httpv1.AppHandler) appErrors.AppHandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) error {
		path := strings.TrimPrefix(req.URL.Path, HttpV1Prefix+"/film")

		adminMiddleware := middleware.AuthRoleMiddleware(constants.AdminRole)
		switch {
		case path == "/search" && http.MethodGet == req.Method:
			return appHandler.FilmHandler.SearchByNameAndActorName(res, req)
		case http.MethodGet == req.Method:
			return appHandler.FilmHandler.GetByQuery(res, req)
		case http.MethodPost == req.Method:
			return adminMiddleware(appHandler.FilmHandler.Create)(res, req)
		case http.MethodPut == req.Method:
			return adminMiddleware(appHandler.FilmHandler.Update)(res, req)
		case http.MethodDelete == req.Method:
			return adminMiddleware(appHandler.FilmHandler.Delete)(res, req)
		default:
			http.NotFound(res, req)
		}
		return nil
	}
}
