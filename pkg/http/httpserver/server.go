package httpserver

import (
	"errors"
	"fmt"
	"net/http"
)

type AppServer struct {
	Server *http.Server
	config *ServerConfig
}

func New(config *ServerConfig) *AppServer {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: nil}

	return &AppServer{
		Server: httpServer,
		config: config,
	}
}

func (server *AppServer) StartAsync() <-chan error {
	exitChannel := make(chan error)

	go func() {
		err := server.Server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			exitChannel <- err
			return
		} else {
			exitChannel <- nil
		}
	}()

	return exitChannel
}
