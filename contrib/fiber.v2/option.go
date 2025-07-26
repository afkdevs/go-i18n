package fiber

import "github.com/gofiber/fiber/v2"

type config struct {
	headerKey   string
	langHandler func(c *fiber.Ctx) string
}

const defaultHeaderKey = "Accept-Language"

// Option is a function that configures the Fiber middleware.
type Option func(*config)

func newConfig(opts ...Option) *config {
	c := &config{
		headerKey: defaultHeaderKey,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithLanguageHandler sets the language handler for the Fiber middleware.
func WithLanguageHandler(handler func(c *fiber.Ctx) string) Option {
	return func(c *config) {
		c.langHandler = handler
	}
}

// WithHeaderKey sets the header key for the Fiber middleware.
//
// Note: It will be ignored if option WithLanguageHandler is set.
func WithHeaderKey(key string) Option {
	return func(c *config) {
		if key == "" {
			return
		}
		c.headerKey = key
	}
}
