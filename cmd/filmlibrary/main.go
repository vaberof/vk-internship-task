package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/vaberof/vk-internship-task/internal/app/entrypoint/http"
	"github.com/vaberof/vk-internship-task/internal/domain"
	pgstorage "github.com/vaberof/vk-internship-task/internal/infra/storage/postgres"
	"github.com/vaberof/vk-internship-task/internal/infra/storage/postgres/pguser"
	"github.com/vaberof/vk-internship-task/internal/service/auth"
	"github.com/vaberof/vk-internship-task/internal/service/user"
	"github.com/vaberof/vk-internship-task/pkg/database/postgres"
	"github.com/vaberof/vk-internship-task/pkg/http/httpserver"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "List of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()
	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	appConfig := mustGetAppConfig(*appConfigPaths)

	fmt.Printf("%+v\n", appConfig)

	logger := logs.New(os.Stdout, nil)

	postgresManagedDb, err := postgres.New(&appConfig.Postgres)
	if err != nil {
		panic(err)
	}

	actorStorage := pgstorage.NewPgActorStorage(postgresManagedDb.PostgresDb)
	filmStorage := pgstorage.NewPgFilmStorage(postgresManagedDb.PostgresDb)
	userStorage := pguser.NewPgUserStorage(postgresManagedDb.PostgresDb)

	actorService := domain.NewActorService(actorStorage, logger)
	filmService := domain.NewFilmService(filmStorage, actorStorage, logger)

	userService := user.NewUserService(userStorage, logger)
	authService := auth.NewAuthService(userService, logger)

	httpRequestBodyValidator := validator.New()

	httpHandler := http.NewHandler(actorService, filmService, authService, httpRequestBodyValidator)

	appServer := httpserver.New(&appConfig.Server, logger)

	httpHandler.InitRoutes(appServer.Mux)

	serverExitChannel := appServer.StartAsync()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case signalValue := <-quitCh:
		logger.GetLogger().Info("stopping application", "signal", signalValue.String())

		gracefulShutdown(appServer, postgresManagedDb)
	case err := <-serverExitChannel:
		logger.GetLogger().Info("stopping application", "err", err.Error())

		gracefulShutdown(appServer, postgresManagedDb)
	}
}

func gracefulShutdown(server *httpserver.AppServer, postgresManagedDb *postgres.ManagedDatabase) {
	if err := server.Server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v\n", err)
	}

	if err := postgresManagedDb.Disconnect(); err != nil {
		log.Printf("Postgres database Shutdown: %v\n", err)
	}

	log.Println("Server successfully shutdown")
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
