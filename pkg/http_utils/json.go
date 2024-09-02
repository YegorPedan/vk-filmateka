package httpUtils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	appValidator "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_validator"
)

const (
	UnmarshalError           = "error unmarshal request body %v"
	ReadBodyError            = "error read request body %v"
	ValidateRequestBodyError = "error validate request body %v"
)

func BodyJson(req *http.Request, body interface{}) error {
	byteBody, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf(ReadBodyError, err)
	}

	if err = json.Unmarshal(byteBody, body); err != nil {
		return fmt.Errorf(UnmarshalError, err)
	}

	validate := appValidator.New()
	if err = validate.Struct(body); err != nil {
		return fmt.Errorf(ValidateRequestBodyError, err)
	}

	return nil
}

func SendJson(res http.ResponseWriter, statusCode int, data interface{}) {
	res.WriteHeader(statusCode)
	encoder := json.NewEncoder(res)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		http.Error(res, "error", http.StatusInternalServerError)
	}
}
