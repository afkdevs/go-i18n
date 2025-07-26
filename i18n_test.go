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

func TestI18n(t *testing.T) {
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
			messageID:       "hello",
			options:         []any{i18n.Param("name", "John")},
			expectedMessage: "Hello, John!",
		},
		{
			name:            "with params",
			messageID:       "hello_age",
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
			messageID:       "only_in_en",
			options:         []any{i18n.Lang("id")},
			expectedMessage: "This message is only available in English.",
		},
		{
			name:            "not found",
			messageID:       "not_found",
			expectedMessage: "ERROR: missing translation for \"not_found\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			message := i18n.T(tc.messageID, tc.options...)
			assert.Equal(t, tc.expectedMessage, message)
		})
	}
}

func TestI18nCtx(t *testing.T) {
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
			messageID:       "hello",
			options:         []any{i18n.Param("name", "John")},
			expectedMessage: "Hello, John!",
		},
		{
			name:            "with params",
			messageID:       "hello_age",
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
			messageID:       "only_in_en",
			language:        "id",
			options:         []any{i18n.Lang("es")},
			expectedMessage: "This message is only available in English.",
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
				ctx = i18n.NewContextWithLanguage(ctx, tc.language)
			}
			message := i18n.TCtx(ctx, tc.messageID, tc.options...)
			assert.Equal(t, tc.expectedMessage, message)
		})
	}
}

func TestI18nWhenTranslationNotFound(t *testing.T) {
	err := i18n.Init(language.English, i18n.WithTranslationFile("testdata/es.yaml"))
	assert.Error(t, err)

	err = i18n.Init(language.English, i18n.WithTranslationFSFile(testdata.FS, "es.yaml"))
	assert.Error(t, err)
}
