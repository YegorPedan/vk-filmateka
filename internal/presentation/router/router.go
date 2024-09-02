package router

import (
	"github.com/OddEer0/vk-filmoteka/internal/presentation/handlers/httpv1"
	"log/slog"
	"net/http"
	"strings"

	appErrors "github.com/OddEer0/vk-filmoteka/internal/common/lib/app_errors"
)

func NewAppRouter(log *slog.Logger, appHandler *httpv1.AppHandler) *http.ServeMux {
	mux := http.NewServeMux()
	middleware := appErrors.LoggingMiddleware(log)

	mux.HandleFunc("/", middleware(func(res http.ResponseWriter, req *http.Request) error {
		switch {
		case strings.HasPrefix(req.URL.Path, HttpV1Prefix):
			return HttpV1Router(appHandler)(res, req)
		default:
			http.NotFound(res, req)
		}
		return nil
	}))

	return mux
}
