package model

type Token struct {
	Id    string `json:"id" validate:"required,uuidv4"`
	Value string `json:"value" validate:"required"`
}
