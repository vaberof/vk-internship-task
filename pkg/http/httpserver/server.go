package httpserver

import (
	"errors"
	"fmt"
	"github.com/vaberof/vk-internship-task/pkg/http/httpserver/middleware/logging"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"log/slog"
	"net/http"
)

type AppServer struct {
	Server *http.Server
	Mux    *http.ServeMux
	config *ServerConfig
	logger *slog.Logger
}

func New(config *ServerConfig, logsBuilder *logs.Logs) *AppServer {
	loggingMw := logging.New(logsBuilder)
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: loggingMw.Handler(mux),
	}

	return &AppServer{
		Server: httpServer,
		Mux:    mux,
		config: config,
		logger: loggingMw.Logger,
	}
}

func (server *AppServer) StartAsync() <-chan error {
	exitChannel := make(chan error)

	server.logger.Info("Starting HTTP server")

	go func() {
		err := server.Server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			server.logger.Error("Failed to start HTTP server")
			exitChannel <- err
			return
		} else {
			exitChannel <- nil
		}
	}()

	server.logger.Info(fmt.Sprintf("Started HTTP server at %s", server.Server.Addr))

	return exitChannel
}
