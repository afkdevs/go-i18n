package fiber_test

import (
	"net/http/httptest"
	"testing"

	"github.com/afkdevs/go-i18n"
	fiberi18n "github.com/afkdevs/go-i18n/contrib/fiber.v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func TestNew(t *testing.T) {
	err := i18n.Init(language.English,
		i18n.WithTranslationFile("../../testdata/en.yaml", "../../testdata/id.yaml"),
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
	)
	require.NoError(t, err, "Failed to initialize i18n")

	app := fiber.New()
	app.Use(fiberi18n.New())

	testHandler := func(c *fiber.Ctx) error {
		return c.SendString(fiberi18n.TCtx(c, "test"))
	}
	helloHandler := func(c *fiber.Ctx) error {
		name := c.Query("name")
		if name == "" {
			return c.SendString(fiberi18n.GetCtx(c, "hello_name"))
		}
		return c.SendString(fiberi18n.GetCtx(c, "hello_name", map[string]any{"name": name}))
	}
	missingHandler := func(c *fiber.Ctx) error {
		return c.SendString(fiberi18n.TCtx(c, "not_exist"))
	}
	app.Get("/test", testHandler)
	app.Get("/hello", helloHandler)
	app.Get("/missing", missingHandler)

	testCases := []struct {
		name       string
		path       string
		acceptLang string
		expected   string
	}{
		{
			name:       "when request has header Accept-Language with id-ID",
			path:       "/test",
			acceptLang: "id-ID",
			expected:   "Ini adalah pesan tes",
		},
		{
			name:       "when request has header Accept-Language with en-US",
			path:       "/test",
			acceptLang: "en-US",
			expected:   "This is test message",
		},
		{
			name:       "when request not has header Accept-Language",
			path:       "/test",
			acceptLang: "",
			expected:   "This is test message",
		},
		{
			name:       "when request has header Accept-Language with id-ID and query name",
			path:       "/hello?name=John",
			acceptLang: "id-ID",
			expected:   "Halo, John",
		},
		{
			name:       "when request has header Accept-Language with en-US and query name",
			path:       "/hello?name=John",
			acceptLang: "en-US",
			expected:   "Hello, John",
		},
		{
			name:       "when request not has header Accept-Language and query name",
			path:       "/hello?name=John",
			acceptLang: "",
			expected:   "Hello, John",
		},
		{
			name:       "when request has header Accept-Language with id-ID and query name is empty",
			path:       "/hello",
			acceptLang: "id-ID",
			expected:   "Halo, <no value>",
		},
		{
			name:     "when translation not found",
			path:     "/missing",
			expected: "ERROR: missing translation for \"not_exist\"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			if tc.acceptLang != "" {
				req.Header.Set("Accept-Language", tc.acceptLang)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err, "Failed to test request")
			assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

			buf := make([]byte, resp.ContentLength)
			resp.Body.Read(buf)
			assert.Equal(t, tc.expected, string(buf), "Expected response body to match")
		})
	}
}
