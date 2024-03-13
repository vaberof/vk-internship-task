package main

import (
	"errors"
	"github.com/vaberof/vk-internship-task/pkg/config"
	"github.com/vaberof/vk-internship-task/pkg/database/postgres"
	"github.com/vaberof/vk-internship-task/pkg/http/httpserver"
	"os"
)

type AppConfig struct {
	Server   httpserver.ServerConfig
	Postgres postgres.Config
}

func mustGetAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var serverConfig httpserver.ServerConfig
	err := config.ParseConfig(provider, "app.http.server", &serverConfig)
	if err != nil {
		return nil, err
	}

	var postgresConfig postgres.Config
	err = config.ParseConfig(provider, "app.postgres", &postgresConfig)
	if err != nil {
		return nil, err
	}
	postgresConfig.User = os.Getenv("POSTGRES_USER")
	postgresConfig.Password = os.Getenv("POSTGRES_PASSWORD")

	appConfig := AppConfig{
		Server:   serverConfig,
		Postgres: postgresConfig,
	}

	return &appConfig, nil
}
