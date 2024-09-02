package main

import (
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/storage/postgres"
	"github.com/OddEer0/vk-filmoteka/internal/presentation/handlers/httpv1"
	"log"
	"net/http"

	_ "github.com/OddEer0/vk-filmoteka/docs"
	"github.com/OddEer0/vk-filmoteka/internal/infrastructure/config"
	slogger "github.com/OddEer0/vk-filmoteka/internal/infrastructure/logger"
	appRouter "github.com/OddEer0/vk-filmoteka/internal/presentation/router"
)

// @title VK-Filmoteka
// @version 1.0
// @description This is a sample HTTP package with Swagger annotations.
func main() {
	cfg := config.MustLoad()
	db, err := postgres.ConnectPg(cfg)
	if err != nil {
		log.Fatal("Error connect postgres", err.Error())
	}
	appHandler := httpv1.NewAppHandler(db)
	logger := slogger.SetupLogger(cfg.Env)
	logger.Info("Logger setup")
	router := appRouter.NewAppRouter(logger, appHandler)
	logger.Info("router setup")
	initSwagger(router)
	logger.Info("swagger setup")

	server := http.Server{Addr: cfg.Server.Address, Handler: router}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server")
	}
}

func initSwagger(mux *http.ServeMux) {
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui/"))))
	mux.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
}
