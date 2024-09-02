package aggregate_test

import (
	"errors"
	"testing"

	"github.com/OddEer0/vk-filmoteka/internal/domain/aggregate"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestActorAggregate(t *testing.T) {
	memData := getActorTestData()

	testCases := []struct {
		name              string
		actorModel        model.Actor
		expectedAggregate *aggregate.ActorAggregate
		isError           bool
		errFields         []string
	}{
		{
			name:              "Should correct create aggregate",
			actorModel:        memData.correctActor,
			expectedAggregate: &aggregate.ActorAggregate{Actor: memData.correctActor},
			isError:           false,
		},
		{
			name:              "Should correct create aggregate 2",
			actorModel:        memData.correctActor2,
			expectedAggregate: &aggregate.ActorAggregate{Actor: memData.correctActor2},
			isError:           false,
		},
		{
			name:              "Should require errors",
			actorModel:        model.Actor{},
			expectedAggregate: &aggregate.ActorAggregate{Actor: model.Actor{}},
			isError:           true,
			errFields:         []string{"Id", "Name", "Gender", "Birthday"},
		},
		{
			name:              "Should id error",
			actorModel:        memData.incorrectIdActor,
			expectedAggregate: &aggregate.ActorAggregate{Actor: memData.incorrectIdActor},
			isError:           true,
			errFields:         []string{"Id"},
		},
		{
			name:              "Should id error",
			actorModel:        memData.incorrectActorGender,
			expectedAggregate: &aggregate.ActorAggregate{Actor: memData.incorrectActorGender},
			isError:           true,
			errFields:         []string{"Gender"},
		},
		{
			name:              "Should birthday error",
			actorModel:        memData.incorrectBirthdayActor,
			expectedAggregate: &aggregate.ActorAggregate{Actor: memData.incorrectBirthdayActor},
			isError:           true,
			errFields:         []string{"Birthday"},
		},
		{
			name:              "Should min error",
			actorModel:        memData.incorrectMinActor,
			expectedAggregate: &aggregate.ActorAggregate{Actor: memData.incorrectMinActor},
			isError:           true,
			errFields:         []string{"Name"},
		},
		{
			name:              "Should max error",
			actorModel:        memData.incorrectMaxActor,
			expectedAggregate: &aggregate.ActorAggregate{Actor: memData.incorrectMaxActor},
			isError:           true,
			errFields:         []string{"Name"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := aggregate.NewActorAggregate(tc.actorModel)
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
