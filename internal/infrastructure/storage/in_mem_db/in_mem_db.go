package inMemDb

import (
	"github.com/OddEer0/vk-filmoteka/internal/common/constants"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/valuesobject"
)

type ActorFilm struct {
	ActorId string
	FilmId  string
}

type InMemDb struct {
	Users     []*model.User
	Tokens    []*model.Token
	Actor     []*model.Actor
	Film      []*model.Film
	ActorFilm []*ActorFilm
}

func (i *InMemDb) CleanUp() {
	i.Users = []*model.User{}
	i.Tokens = []*model.Token{}
	i.Actor = []*model.Actor{}
	i.Film = []*model.Film{}
	i.ActorFilm = []*ActorFilm{}
}

var instance *InMemDb = nil

func New() *InMemDb {
	if instance != nil {
		return instance
	}

	instance = &InMemDb{
		Users:     []*model.User{},
		Tokens:    []*model.Token{},
		Actor:     []*model.Actor{},
		Film:      []*model.Film{},
		ActorFilm: []*ActorFilm{},
	}

	password, _ := valuesobject.NewPassword("Adminadmin41")

	instance.Users = append(instance.Users, &model.User{
		Id:       "admin",
		Name:     "Admin",
		Password: password,
		Role:     constants.AdminRole,
	})

	return instance
}
