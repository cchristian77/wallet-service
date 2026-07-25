package fhttp

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
