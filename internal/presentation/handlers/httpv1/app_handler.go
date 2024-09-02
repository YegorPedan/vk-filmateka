package httpv1

import (
	"database/sql"
	tokenService "github.com/OddEer0/vk-filmoteka/internal/app/services/token_service"
	userService "github.com/OddEer0/vk-filmoteka/internal/app/services/user_service"
	actorUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/actor_usecase"
	authUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/auth_usecase"
	filmUseCase "github.com/OddEer0/vk-filmoteka/internal/app/usecases/film_usecase"
	mockRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/mock_repository"
	postgresRepository "github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/postgres_repository"
)

type (
	AppHandler struct {
		AuthHandler
		FilmHandler
		ActorHandler
	}
)

var instance *AppHandler = nil
var instance2 *AppHandler = nil

func NewAppHandler(db *sql.DB) *AppHandler {
	if instance != nil {
		return instance
	}

	userRepo := postgresRepository.NewUserRepository(db)
	tokenRepo := postgresRepository.NewTokenRepository(db)
	actorRepo := postgresRepository.NewActorRepository(db)
	filmRepo := postgresRepository.NewFilmRepository(db)

	userServ := userService.New(userRepo)
	tokenServ := tokenService.New(tokenRepo)

	authUsecase := authUseCase.New(userServ, tokenServ, userRepo)
	actorUsecase := actorUseCase.New(actorRepo, filmRepo)
	filmUsecase := filmUseCase.New(filmRepo)

	instance = &AppHandler{
		AuthHandler:  NewAuthHandler(authUsecase),
		FilmHandler:  NewFilmHandler(filmUsecase),
		ActorHandler: NewActorHandler(actorUsecase),
	}

	return instance
}

func NewAppHandlerMock() *AppHandler {
	if instance2 != nil {
		return instance2
	}

	userRepo := mockRepository.NewUserRepository()
	tokenRepo := mockRepository.NewTokenRepository()
	actorRepo := mockRepository.NewActorRepository()
	filmRepo := mockRepository.NewFilmRepository()

	userServ := userService.New(userRepo)
	tokenServ := tokenService.New(tokenRepo)

	authUsecase := authUseCase.New(userServ, tokenServ, userRepo)
	actorUsecase := actorUseCase.New(actorRepo, filmRepo)
	filmUsecase := filmUseCase.New(filmRepo)

	instance2 = &AppHandler{
		AuthHandler:  NewAuthHandler(authUsecase),
		FilmHandler:  NewFilmHandler(filmUsecase),
		ActorHandler: NewActorHandler(actorUsecase),
	}

	return instance2
}
