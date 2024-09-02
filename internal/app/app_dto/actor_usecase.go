package appDto

import (
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"time"
)

type (
	CreateActorUseCaseDto struct {
		Name     string    `json:"name" validate:"required,min=1,max=100"`
		Gender   string    `json:"gender" validate:"required,gender"`
		Birthday time.Time `json:"birthday" validate:"required,dateIsLessNow"`
	}

	ActorGetByQueryResult struct {
		Actors    []*aggregate.ActorAggregate `json:"actors"`
		PageCount int                         `json:"pageCount"`
	}
)
