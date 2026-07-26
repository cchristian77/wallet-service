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

// NewHTTPServer creates an HTTP server for the given handler e.g., ServeMux with middlewares.
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

// Start - Starts the http server
func (s *Server) Start(ctx context.Context) (err error) {
	logger.L().Info(fmt.Sprintf("Starting Http Server at {%v}....", s.server.Addr))

	timeout := time.Duration(s.options.serverStopTimeoutSeconds) * time.Second

	idleConnClosed := make(chan struct{})

	// For graceful shutdown
	go func() {
		signal.Notify(s.exit, os.Interrupt, syscall.SIGTERM)

		sig := <-s.exit
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		logger.L().Info(fmt.Sprintf("Received signal {%v}, shutting down HTTP Server....", sig.String()))

		if err = s.server.Shutdown(ctx); err != nil {
			logger.L().Error(fmt.Sprintf("Error while shutting down HTTP server: %v", err))
		}

		close(idleConnClosed)
	}()

	// Start the server
	// here, we use a different var `serveErr` as ListenAndServe() always returns a non-nil error, and we want to
	// return the error only if it's not http.ErrServerClosed
	if serveErr := s.server.ListenAndServe(); !errors.Is(serveErr, http.ErrServerClosed) {
		logger.L().Error(fmt.Sprintf("Error starting up Http Server. Error: %v", serveErr.Error()))
		err = serveErr
	}

	// Wait for the shutdown process to complete
	<-idleConnClosed

	return err
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) {
	logger.L().Info("Shutting down Http Server....")
	s.exit <- os.Interrupt
}
