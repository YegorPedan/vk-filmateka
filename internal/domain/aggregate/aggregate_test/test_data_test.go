package aggregate_test

import (
	"strings"
	"time"

	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/OddEer0/vk-filmoteka/internal/domain/valuesobject"
	"github.com/google/uuid"
)

type InMemUserTestData struct {
	correctUser       model.User
	correctUser2      model.User
	incorrectIdUser   model.User
	incorrectUserRole model.User
	incorrectMinUser  model.User
	incorrectMaxUser  model.User
}

type InMemActorTestData struct {
	correctActor           model.Actor
	correctActor2          model.Actor
	incorrectIdActor       model.Actor
	incorrectActorGender   model.Actor
	incorrectMinActor      model.Actor
	incorrectMaxActor      model.Actor
	incorrectBirthdayActor model.Actor
}

type InMemFilmTestData struct {
	correctFilm      model.Film
	correctFilm2     model.Film
	incorrectIdFilm  model.Film
	incorrectMinFilm model.Film
	incorrectMaxFilm model.Film
}

func getUserTestData() *InMemUserTestData {
	correctUserPassword, _ := valuesobject.NewPassword("correct12")

	return &InMemUserTestData{
		correctUser: model.User{
			Id:       uuid.New().String(),
			Name:     "Marlen",
			Password: correctUserPassword,
			Role:     "ADMIN",
		},
		correctUser2: model.User{
			Id:       uuid.New().String(),
			Name:     "Marlen",
			Password: correctUserPassword,
			Role:     "USER",
		},
		incorrectIdUser: model.User{
			Id:       "notuuidv4",
			Name:     "Marlen",
			Password: correctUserPassword,
			Role:     "ADMIN",
		},
		incorrectUserRole: model.User{
			Id:       uuid.New().String(),
			Name:     "Marlen",
			Password: correctUserPassword,
			Role:     "NOADMIN",
		},
		incorrectMinUser: model.User{
			Id:       uuid.New().String(),
			Name:     "Ma",
			Password: correctUserPassword,
			Role:     "ADMIN",
		},
		incorrectMaxUser: model.User{
			Id:       uuid.New().String(),
			Name:     strings.Repeat("longlong12", 100) + "1",
			Password: correctUserPassword,
			Role:     "ADMIN",
		},
	}
}

func getActorTestData() *InMemActorTestData {
	return &InMemActorTestData{
		correctActor:           model.Actor{Id: uuid.New().String(), Name: "Jason", Gender: "male", Birthday: time.Now().AddDate(-20, 0, 0)},
		correctActor2:          model.Actor{Id: uuid.New().String(), Name: "Katya", Gender: "female", Birthday: time.Now().AddDate(-25, 0, 0)},
		incorrectIdActor:       model.Actor{Id: "NOTUUIDV4", Name: "Jason", Gender: "male", Birthday: time.Now().AddDate(-20, 0, 0)},
		incorrectActorGender:   model.Actor{Id: uuid.New().String(), Name: "Katya", Gender: "trans", Birthday: time.Now().AddDate(-25, 0, 0)},
		incorrectBirthdayActor: model.Actor{Id: uuid.New().String(), Name: "Jason", Gender: "male", Birthday: time.Now().AddDate(0, 0, 1)},
		incorrectMinActor:      model.Actor{Id: uuid.New().String(), Name: "", Gender: "male", Birthday: time.Now().AddDate(-20, 0, 0)},
		incorrectMaxActor:      model.Actor{Id: uuid.New().String(), Name: strings.Repeat("longlong21", 10) + "1", Gender: "male", Birthday: time.Now().AddDate(-20, 0, 0)},
	}
}

func getFilmTestData() *InMemFilmTestData {
	description2 := "cool film"
	incorrectMaxDescription := strings.Repeat("longlong34", 100) + "1"

	return &InMemFilmTestData{
		correctFilm:      model.Film{Id: uuid.New().String(), Name: "Titanic", Description: nil, ReleaseDate: time.Now().AddDate(-17, 0, 0), Rate: 10},
		correctFilm2:     model.Film{Id: uuid.New().String(), Name: "Titanic", Description: &description2, ReleaseDate: time.Now().AddDate(-17, 0, 0), Rate: 0},
		incorrectIdFilm:  model.Film{Id: "NOTUUIDV4", Name: "Titanic", Description: nil, ReleaseDate: time.Now().AddDate(-17, 0, 0), Rate: 10},
		incorrectMaxFilm: model.Film{Id: uuid.New().String(), Name: strings.Repeat("longlong32", 15) + "1", Description: &incorrectMaxDescription, ReleaseDate: time.Now().AddDate(-17, 0, 0), Rate: 11},
		incorrectMinFilm: model.Film{Id: uuid.New().String(), Name: "", Description: nil, ReleaseDate: time.Now().AddDate(-17, 0, 0), Rate: -1},
	}
}
