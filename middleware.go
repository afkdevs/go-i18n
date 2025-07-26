package i18n

import (
	"context"
	"net/http"

	"golang.org/x/text/language"
)

type contextKey string

const languageCtxKey contextKey = "i18n-language"

func defaultLanguageHandler(headerKey string) func(r *http.Request) string {
	return func(r *http.Request) string {
		return r.Header.Get(headerKey)
	}
}

// NewMiddleware creates a middleware that sets the language to the context from the request.
//
// It uses the Accept-Language header to get the language.
func NewMiddleware(opts ...MiddlewareOption) func(http.Handler) http.Handler {
	cfg := newMiddlewareConfig(opts...)
	if cfg.langHandler == nil {
		cfg.langHandler = defaultLanguageHandler(cfg.headerKey)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := cfg.langHandler(r)
			if lang != "" {
				ctx := NewContextWithLanguage(r.Context(), lang)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetLanguage returns the language tag from the context.
//
// If the language tag is not found, it returns the default language tag.
func GetLanguage(ctx context.Context) language.Tag {
	lang, ok := ctx.Value(languageCtxKey).(string)
	if !ok {
		return defaultLanguage
	}
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil || len(tags) == 0 {
		return defaultLanguage
	}
	return tags[0]
}

// NewContextWithLanguage sets the language to the context.
//
// You can use this function to set the language to the context manually.
func NewContextWithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, languageCtxKey, language)
}
