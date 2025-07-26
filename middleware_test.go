package i18n_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afkdevs/go-i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func chainMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}

func TestMiddleware(t *testing.T) {
	err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("testdata/en.yaml", "testdata/id.yaml"),
	)
	require.NoError(t, err)

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		message := i18n.TCtx(r.Context(), "test")
		_, _ = w.Write([]byte(message))
	}
	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		message := i18n.TCtx(r.Context(), "hello_name", i18n.Param("name", name))
		_, _ = w.Write([]byte(message))
	}
	i18nMiddleware := i18n.NewMiddleware(
		i18n.WithHeaderKey("Accept-Language"),
	)

	r := http.NewServeMux()
	r.Handle("/test", chainMiddleware(http.HandlerFunc(testHandler), i18nMiddleware))
	r.Handle("/hello", chainMiddleware(http.HandlerFunc(helloHandler), i18nMiddleware))

	testCases := []struct {
		name            string
		acceptLanguage  string
		path            string
		expectedMessage string
	}{
		{
			name:            "without header and query param",
			path:            "/test",
			expectedMessage: "This is test message",
		},
		{
			name:            "without header and with query param",
			path:            "/hello?name=John",
			expectedMessage: "Hello, John",
		},
		{
			name:            "with accept-language en",
			path:            "/test",
			acceptLanguage:  "en",
			expectedMessage: "This is test message",
		},
		{
			name:            "with accept-language en and with query param",
			path:            "/hello?name=John",
			acceptLanguage:  "en",
			expectedMessage: "Hello, John",
		},
		{
			name:            "with accept-language id",
			path:            "/test",
			acceptLanguage:  "id",
			expectedMessage: "Ini adalah pesan tes",
		},
		{
			name:            "with accept-language id and with query param",
			path:            "/hello?name=John",
			acceptLanguage:  "id",
			expectedMessage: "Halo, John",
		},
		{
			name:            "with accept-language es",
			path:            "/test",
			acceptLanguage:  "es",
			expectedMessage: "This is test message",
		},
		{
			name:            "with multiple accept-language",
			path:            "/test",
			acceptLanguage:  "es-ES,id-ID,en-US",
			expectedMessage: "Ini adalah pesan tes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			if tc.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tc.acceptLanguage)
			}
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, tc.expectedMessage, resp.Body.String())
		})
	}

	r = http.NewServeMux()

	i18nMiddleware = i18n.NewMiddleware(
		i18n.WithLanguageHandler(func(r *http.Request) string {
			return r.URL.Query().Get("lang")
		}),
	)
	r.Handle("/test", chainMiddleware(http.HandlerFunc(testHandler), i18nMiddleware))

	testCasesCustomHandler := []struct {
		name            string
		path            string
		expectedMessage string
	}{
		{
			name:            "with query lang en",
			path:            "/test?lang=en",
			expectedMessage: "This is test message",
		},
		{
			name:            "with query lang id",
			path:            "/test?lang=id",
			expectedMessage: "Ini adalah pesan tes",
		},
		{
			name:            "with query lang empty",
			path:            "/test",
			expectedMessage: "This is test message",
		},
	}
	for _, tc := range testCasesCustomHandler {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)
			assert.Equal(t, tc.expectedMessage, resp.Body.String())
		})
	}
}

func TestGetLanguage(t *testing.T) {
	err := i18n.Init(language.English)
	require.NoError(t, err)

	testCases := []struct {
		name string
		lang string
		tag  language.Tag
	}{
		{
			name: "blank",
			tag:  language.English,
		},
		{
			name: "english language",
			lang: "en",
			tag:  language.English,
		},
		{
			name: "indonesian language",
			lang: "id",
			tag:  language.Indonesian,
		},
		{
			name: "multiple language",
			lang: "id,en",
			tag:  language.Indonesian,
		},
		{
			name: "invalid language",
			lang: "invalid",
			tag:  language.English,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.lang != "" {
				ctx = i18n.SetLangToContext(ctx, tc.lang)
			}
			tag := i18n.GetLanguage(ctx)
			assert.Equal(t, tc.tag, tag)
		})
	}
}
