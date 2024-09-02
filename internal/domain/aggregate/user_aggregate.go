package aggregate

import (
	appValidator "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_validator"
	"github.com/OddEer0/vk-filmoteka/internal/domain/model"
)

type UserAggregate struct {
	User  model.User
	Token *model.Token
}

func (u *UserAggregate) Validation() error {
	validator := appValidator.New()
	err := validator.Struct(u.User)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserAggregate) SetToken(refreshToken string) error {
	token := model.Token{Id: u.User.Id, Value: refreshToken}
	validator := appValidator.New()
	err := validator.Struct(token)
	if err != nil {
		return err
	}
	u.Token = &token
	return nil
}

func NewUserAggregate(user model.User) (*UserAggregate, error) {
	result := &UserAggregate{User: user}
	if err := result.Validation(); err != nil {
		return nil, err
	}
	return result, nil
}
