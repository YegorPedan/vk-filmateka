package appValidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func uuidv4(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	parsedUUID, err := uuid.Parse(value)
	return err == nil && parsedUUID.Version() == uuid.Version(4)
}
