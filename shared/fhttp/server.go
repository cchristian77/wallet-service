package fhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cchristian77/wallet-service/util/config"
	"github.com/cchristian77/wallet-service/util/logger"
)

// Server wraps http.Server with graceful shutdown.
type Server struct {
	server  *http.Server
	options *serverOptions
	exit    chan os.Signal
}

// NewHTTPServer creates an HTTP server for the given handler (e.g. ServeMux + middleware).
func NewHTTPServer(handler http.Handler, opts ...Option) (*Server, error) {
	options := &serverOptions{
		serverStopTimeoutSeconds: defaultStopTimeoutSeconds,
		readTimeoutSeconds:       defaultReadTimeoutSeconds,
		readHeaderTimeoutSeconds: defaultReadHeaderTimeoutSeconds,
	}

	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}

	if config.Env() == nil {
		return nil, errors.New("config is nil")
	}

	if config.Env().App.Port < 1 {
		return nil, errors.New("port is less than 1")
	}

	if handler == nil {
		return nil, errors.New("handler should not be nil")
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%v", config.Env().App.Port),
		Handler:           handler,
		ReadTimeout:       time.Duration(options.readTimeoutSeconds) * time.Second,
		ReadHeaderTimeout: time.Duration(options.readHeaderTimeoutSeconds) * time.Second,
	}

	return &Server{
		server:  httpServer,
		options: options,
		exit:    make(chan os.Signal, 1),
	}, nil
}

// Start listens until SIGINT/SIGTERM, then shuts down gracefully.
func (s *Server) Start(_ context.Context) error {
	logger.L().Info(fmt.Sprintf("Starting Http Server at {%v}....", s.server.Addr))

	timeout := time.Duration(s.options.serverStopTimeoutSeconds) * time.Second
	idleConnClosed := make(chan struct{})

	go func() {
		signal.Notify(s.exit, os.Interrupt, syscall.SIGTERM)

		sig := <-s.exit
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		logger.L().Info(fmt.Sprintf("Received signal {%v}, shutting down HTTP Server....", sig.String()))

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			logger.L().Error(fmt.Sprintf("Error while shutting down HTTP server: %v", err))
		}

		close(idleConnClosed)
	}()

	var err error
	if serveErr := s.server.ListenAndServe(); !errors.Is(serveErr, http.ErrServerClosed) {
		logger.L().Error(fmt.Sprintf("Error starting up Http Server. Error: %v", serveErr.Error()))
		err = serveErr
	}

	<-idleConnClosed

	return err
}

// Stop triggers graceful shutdown.
func (s *Server) Stop(_ context.Context) {
	logger.L().Info("Shutting down Http Server....")
	s.exit <- os.Interrupt
}
