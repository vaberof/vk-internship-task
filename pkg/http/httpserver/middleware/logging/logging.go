package logging

import (
	"fmt"
	"github.com/vaberof/vk-internship-task/pkg/logging/logs"
	"log/slog"
	"net/http"
)

type Middleware struct {
	Handler func(http.Handler) http.Handler
	Logger  *slog.Logger
}

func New(logs *logs.Logs) *Middleware {
	return impl(logs, "")
}

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriterWrapper) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func impl(logs *logs.Logs, serverName string) *Middleware {
	loggerName := "http-server"
	if serverName != "" {
		loggerName = fmt.Sprintf("%s.%s", loggerName, serverName)
	}
	logger := logs.WithName(loggerName)

	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			path := request.URL.Path
			if path == "" {
				path = "/"
			}
			method := request.Method
			logger.Info("Request started", slog.Group("http", "path", path, "method", method))

			writerWrapper := &responseWriterWrapper{
				ResponseWriter: writer,
			}

			defer func() {
				status := writerWrapper.status

				if status >= 500 {
					logger.Info("Request finished", slog.Group("http", "path", path, "method", method, "result", "error", "status", status))
					return
				}

				logger.Info("Request finished", slog.Group("http", "path", path, "method", method, "result", "success", "status", status))
			}()

			next.ServeHTTP(writerWrapper, request)
		})
	}

	return &Middleware{
		Handler: handler,
		Logger:  logger,
	}
}
