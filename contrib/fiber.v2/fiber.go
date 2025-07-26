package fiber

import (
	"github.com/afkdevs/go-i18n"
	"github.com/gofiber/fiber/v2"
)

func defaultLanguageHandler(headerKey string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		return c.Get(headerKey)
	}
}

// New creates a Fiber middleware that sets the language to the context from the request.
//
// Defaults to using the Accept-Language header to get the language.
// You can customize the header key or the language handler using options.
func New(opts ...Option) fiber.Handler {
	cfg := newConfig(opts...)
	if cfg.headerKey == "" {
		cfg.headerKey = defaultHeaderKey
	}
	if cfg.langHandler == nil {
		cfg.langHandler = defaultLanguageHandler(cfg.headerKey)
	}

	return func(c *fiber.Ctx) error {
		lang := cfg.langHandler(c)
		if lang != "" {
			ctx := i18n.SetLangToContext(c.UserContext(), lang)
			c.SetUserContext(ctx)
		}
		return c.Next()
	}
}

// TCtx is an alias for i18n.TCtx that uses the Fiber context.
func TCtx(c *fiber.Ctx, id string, opts ...any) string {
	ctx := c.UserContext()
	return i18n.TCtx(ctx, id, opts...)
}

// GetCtx is an alias for i18n.GetCtx that uses the Fiber context.
func GetCtx(c *fiber.Ctx, id string, opts ...any) string {
	ctx := c.UserContext()
	return i18n.GetCtx(ctx, id, opts...)
}
