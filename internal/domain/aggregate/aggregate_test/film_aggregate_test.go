package aggregate_test

import (
	"errors"
	"testing"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestFilmAggregate(t *testing.T) {
	memData := getFilmTestData()

	testCases := []struct {
		name              string
		filmModel         model.Film
		expectedAggregate *aggregate.FilmAggregate
		isError           bool
		errFields         []string
	}{
		{
			name:              "Should correct create aggregate",
			filmModel:         memData.correctFilm,
			expectedAggregate: &aggregate.FilmAggregate{Film: memData.correctFilm},
			isError:           false,
		},
		{
			name:              "Should correct create aggregate 2",
			filmModel:         memData.correctFilm2,
			expectedAggregate: &aggregate.FilmAggregate{Film: memData.correctFilm2},
			isError:           false,
		},
		{
			name:              "Should required errors",
			filmModel:         model.Film{},
			expectedAggregate: &aggregate.FilmAggregate{Film: model.Film{}},
			isError:           true,
			errFields:         []string{"Id", "Name", "ReleaseDate", "Rate"},
		},
		{
			name:              "Should uuidv4 error",
			filmModel:         memData.incorrectIdFilm,
			expectedAggregate: &aggregate.FilmAggregate{Film: memData.incorrectIdFilm},
			isError:           true,
			errFields:         []string{"Id"},
		},
		{
			name:              "Should min error",
			filmModel:         memData.incorrectMinFilm,
			expectedAggregate: &aggregate.FilmAggregate{Film: memData.incorrectMinFilm},
			isError:           true,
			errFields:         []string{"Name", "Rate"},
		},
		{
			name:              "Should max error",
			filmModel:         memData.incorrectMaxFilm,
			expectedAggregate: &aggregate.FilmAggregate{Film: memData.incorrectMaxFilm},
			isError:           true,
			errFields:         []string{"Name", "Description", "Rate"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := aggregate.NewFilmAggregate(tc.filmModel)
			if tc.isError {
				assert.Error(t, err)
				var validationError validator.ValidationErrors
				ok := errors.As(err, &validationError)
				assert.True(t, ok)
				for i, e := range validationError {
					assert.Equal(t, tc.errFields[i], e.Field())
				}
			} else {
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedAggregate, result)
			}
		})
	}
}
