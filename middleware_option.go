package i18n

import "net/http"

type middlewareConfig struct {
	headerKey   string
	langHandler func(r *http.Request) string
}

const defaultHeaderKey = "Accept-Language"

// MiddlewareOption is a function that configures the middleware.
type MiddlewareOption func(*middlewareConfig)

func newMiddlewareConfig(opts ...MiddlewareOption) *middlewareConfig {
	cfg := &middlewareConfig{
		headerKey: defaultHeaderKey,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithHeaderKey sets the header key for the middleware.
func WithHeaderKey(key string) MiddlewareOption {
	return func(cfg *middlewareConfig) {
		if key == "" {
			return
		}
		cfg.headerKey = key
	}
}

// WithLanguageHandler sets the language handler for the middleware.
func WithLanguageHandler(handler func(r *http.Request) string) MiddlewareOption {
	return func(cfg *middlewareConfig) {
		cfg.langHandler = handler
	}
}
