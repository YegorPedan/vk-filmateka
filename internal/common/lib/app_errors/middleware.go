package appErrors

import (
	"errors"
	"log/slog"
	"net/http"

	httpUtils "github.com/OddEer0/vk-filmoteka/pkg/http_utils"
)

type AppHandlerFunc func(res http.ResponseWriter, req *http.Request) error

func LoggingMiddleware(logger *slog.Logger) func(handlerFunc AppHandlerFunc) http.HandlerFunc {
	return func(handlerFunc AppHandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			err := handlerFunc(res, req)
			if err != nil {
				var appErr *AppError
				if errors.As(err, &appErr) {
					if appErr.Code >= 500 {
						logger.Error("ERROR", "statusCode", appErr.Code, "errorMessage", appErr.Message, "developerMessage", appErr.DevMessage)
					} else if appErr.Code >= 400 {
						logger.Info("INFO", "statusCode", appErr.Code, "errorMessage", appErr.Message, "developerMessage", appErr.DevMessage)
					}
					httpUtils.SendJson(res, appErr.Code, ResponseError{Code: appErr.Code, Message: appErr.Message})
				}
			}
		}
	}
}
