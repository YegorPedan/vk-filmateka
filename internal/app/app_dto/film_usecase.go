package appDto

import (
	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"time"
)

type (
	CreateFilmUseCaseDto struct {
		Name        string    `json:"name" validate:"required,min=1,max=150"`
		Description *string   `json:"description,omitempty" validate:"omitempty,max=1000"`
		ReleaseDate time.Time `json:"release" validate:"required"`
		Rate        float32   `json:"rate" validate:"min=0,max=10"`
	}

	FilmGetByQueryResult struct {
		Films     []*aggregate.FilmAggregate `json:"films"`
		PageCount int                        `json:"pageCount"`
	}
)
