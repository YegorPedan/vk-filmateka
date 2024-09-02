package dto

type (
	RegistrationDto struct {
		Name     string `json:"name" validate:"required,min=3,max=100"`
		Password string `json:"password" validate:"required,min=8,max=35,isValidPassword"`
	}

	LoginDto struct {
		Name     string `json:"name" validate:"required,min=3,max=100"`
		Password string `json:"password" validate:"required,min=8,max=35,isValidPassword"`
	}
)
