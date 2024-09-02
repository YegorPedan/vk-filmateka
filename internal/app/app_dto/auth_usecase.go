package appDto

type (
	RegistrationUseCaseDto struct {
		Name     string `json:"name" validate:"required,min=3,max=100"`
		Password string `json:"password" validate:"required,min=8,max=35,isValidPassword"`
	}

	LoginUseCaseDto struct {
		Name     string `json:"name" validate:"required,min=3,max=100"`
		Password string `json:"password" validate:"required,min=8,max=35,isValidPassword"`
	}

	ResponseUserDto struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}
)
