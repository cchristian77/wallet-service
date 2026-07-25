package fhttp

import (
	"errors"
)

const (
	defaultStopTimeoutSeconds       = 2
	defaultReadTimeoutSeconds       = 10
	defaultReadHeaderTimeoutSeconds = 5
)

type serverOptions struct {
	serverStopTimeoutSeconds int
	readTimeoutSeconds       int
	readHeaderTimeoutSeconds int
}

type Option func(*serverOptions) error

func WithTimeout(timeout int) Option {
	return func(option *serverOptions) error {
		if timeout > 0 {
			option.serverStopTimeoutSeconds = timeout
			return nil
		}

		return errors.New("timeout should be greater than zero")
	}
}

func WithReadTimeout(timeoutSeconds int) Option {
	return func(option *serverOptions) error {
		if timeoutSeconds > 0 {
			option.readTimeoutSeconds = timeoutSeconds
			return nil
		}

		return errors.New("read timeout should be greater than zero")
	}
}

func WithReadHeaderTimeout(timeoutSeconds int) Option {
	return func(option *serverOptions) error {
		if timeoutSeconds > 0 {
			option.readHeaderTimeoutSeconds = timeoutSeconds
			return nil
		}

		return errors.New("read header timeout should be greater than zero")
	}
}
