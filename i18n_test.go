package i18n_test

import (
	"context"
	"testing"

	"github.com/afkdevs/go-i18n"
	"github.com/afkdevs/go-i18n/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func TestT(t *testing.T) {
	t.Run("not initialized", func(t *testing.T) {
		msg := i18n.T("test")
		assert.Equal(t, "ERROR: i18n is not initialized", msg)
	})
	err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("testdata/en.yaml", "testdata/id.yaml"),
	)
	require.NoError(t, err)

	testCases := []struct {
		name            string
		messageID       string
		options         []any
		expectedMessage string
	}{
		{
			name:            "default",
			messageID:       "test",
			expectedMessage: "This is test message",
		},
		{
			name:            "with custom language",
			messageID:       "test",
			options:         []any{i18n.Lang("id")},
			expectedMessage: "Ini adalah pesan tes",
		},
		{
			name:            "with param",
			messageID:       "hello_name",
			options:         []any{i18n.Param("name", "John")},
			expectedMessage: "Hello, John",
		},
		{
			name:            "with params",
			messageID:       "hello_name_age",
			options:         []any{i18n.Params{"name": "John", "age": 30}},
			expectedMessage: "Hello, John! You are 30 years old.",
		},
		{
			name:            "with default message",
			messageID:       "with_default_message",
			options:         []any{i18n.Default("This is default message")},
			expectedMessage: "This is default message",
		},
		{
			name:            "not found and use default language",
			messageID:       "hello_english",
			options:         []any{i18n.Lang("id")},
			expectedMessage: "Hello, This message is only available in English.",
		},
		{
			name:            "not found",
			messageID:       "not_found",
			expectedMessage: "ERROR: missing translation for \"not_found\"",
		},
		{
			name:            "invalid param type",
			messageID:       "hello_name",
			options:         []any{invalidParamFunc()},
			expectedMessage: "Hello, <no value>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			message := i18n.T(tc.messageID, tc.options...)
			assert.Equal(t, tc.expectedMessage, message)
		})
	}
}

func invalidParamFunc() func() string {
	return func() string {
		return "invalid_param"
	}
}

func TestTCtx(t *testing.T) {
	err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFSFile(testdata.FS, "en.yaml", "id.yaml"),
	)
	require.NoError(t, err)

	testCases := []struct {
		name            string
		messageID       string
		options         []any
		expectedMessage string
		language        string
	}{
		{
			name:            "default",
			messageID:       "test",
			expectedMessage: "This is test message",
		},
		{
			name:            "with param",
			messageID:       "hello_name",
			options:         []any{i18n.Param("name", "John")},
			expectedMessage: "Hello, John",
		},
		{
			name:            "with params",
			messageID:       "hello_name_age",
			options:         []any{i18n.Params{"name": "John", "age": 30}},
			expectedMessage: "Hello, John! You are 30 years old.",
		},
		{
			name:            "with default message",
			messageID:       "with_default_message",
			options:         []any{i18n.Default("This is default message")},
			expectedMessage: "This is default message",
		},
		{
			name:            "with custom language",
			messageID:       "test",
			language:        "id",
			expectedMessage: "Ini adalah pesan tes",
		},
		{
			name:            "not found and use default language",
			messageID:       "hello_english",
			language:        "id",
			options:         []any{i18n.Lang("es")},
			expectedMessage: "Hello, This message is only available in English.",
		},
		{
			name:            "not found",
			messageID:       "not_found",
			expectedMessage: "ERROR: missing translation for \"not_found\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.language != "" {
				ctx = i18n.SetLangToContext(ctx, tc.language)
			}
			message := i18n.TCtx(ctx, tc.messageID, tc.options...)
			assert.Equal(t, tc.expectedMessage, message)
		})
	}
}

func TestInit(t *testing.T) {
	t.Run("when translation files not found", func(t *testing.T) {
		err := i18n.Init(language.English, i18n.WithTranslationFile("testdata/es.yaml"))
		assert.Error(t, err)
	})
	t.Run("when translation files not found in FS", func(t *testing.T) {
		err := i18n.Init(language.English, i18n.WithTranslationFSFile(testdata.FS, "es.yaml"))
		assert.Error(t, err)
	})
}
