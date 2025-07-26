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

func TestMiddleware(t *testing.T) {
	err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("testdata/en.yaml", "testdata/id.yaml"),
	)
	require.NoError(t, err)

	handlerA := func(w http.ResponseWriter, r *http.Request) {
		message := i18n.TCtx(r.Context(), "test")
		_, _ = w.Write([]byte(message))
	}
	handlerB := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		message := i18n.TCtx(r.Context(), "hello", i18n.Param("name", name))
		_, _ = w.Write([]byte(message))
	}

	testCases := []struct {
		name            string
		lang            string
		url             string
		handler         http.HandlerFunc
		expectedMessage string
	}{
		{
			name:            "without header and query param",
			url:             "/",
			handler:         handlerA,
			expectedMessage: "This is test message",
		},
		{
			name:            "without header and with query param",
			url:             "/?name=John",
			handler:         handlerB,
			expectedMessage: "Hello, John!",
		},
		{
			name:            "with accept-language en",
			url:             "/",
			lang:            "en",
			handler:         handlerA,
			expectedMessage: "This is test message",
		},
		{
			name:            "with accept-language en and with query param",
			url:             "/?name=John",
			lang:            "en",
			handler:         handlerB,
			expectedMessage: "Hello, John!",
		},
		{
			name:            "with accept-language id",
			url:             "/",
			lang:            "id",
			handler:         handlerA,
			expectedMessage: "Ini adalah pesan tes",
		},
		{
			name:            "with accept-language id and with query param",
			url:             "/?name=John",
			lang:            "id",
			handler:         handlerB,
			expectedMessage: "Halo John",
		},
		{
			name:            "with accept-language es",
			url:             "/",
			lang:            "es",
			handler:         handlerA,
			expectedMessage: "This is test message",
		},
		{
			name:            "with multiple accept-language",
			url:             "/",
			lang:            "es-ES,id-ID,en-US",
			handler:         handlerA,
			expectedMessage: "Ini adalah pesan tes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, tc.url, nil)
			rec := httptest.NewRecorder()
			if tc.lang != "" {
				req.Header.Set("Accept-Language", tc.lang)
			}

			i18n.Middleware(tc.handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedMessage, rec.Body.String())
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
				ctx = i18n.NewContextWithLanguage(ctx, tc.lang)
			}
			tag := i18n.GetLanguage(ctx)
			assert.Equal(t, tc.tag, tag)
		})
	}
}
