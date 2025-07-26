package i18n

import (
	"context"
	"fmt"
	"slices"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle                    *i18n.Bundle
	defaultLanguage           language.Tag
	missingTranslationHandler func(string, error) string
	extractLanguageFunc       func(context.Context) string
)

func defaultExtractLanguageFunc(ctx context.Context) string {
	lang, ok := ctx.Value(languageCtxKey).(string)
	if ok {
		return lang
	}
	return ""
}

func defaultMissingTranslationFunc(messageID string, _ error) string {
	return fmt.Sprintf("ERROR: missing translation for %q", messageID)
}

// Init initializes the i18n package. It must be called before any other function.
//
// Example:
//
//	if err := i18n.Init(language.English,
//		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
//		i18n.WithMessageFilePaths("locales/en.yaml", "locales/id.yaml"),
//	) err != nil {
//		panic(err)
//	}
func Init(language language.Tag, opts ...Option) error {
	defaultOpts := []Option{
		WithMissingTranslationHandler(defaultMissingTranslationFunc),
		WithExtractLanguageFunc(defaultExtractLanguageFunc),
	}
	opts = append(defaultOpts, opts...)
	config := newI18nConfig(opts...)

	defaultLanguage = language
	missingTranslationHandler = config.missingTranslationHandler
	extractLanguageFunc = config.extractLanguageFunc

	bundle = i18n.NewBundle(language)
	for format, unmarshalFunc := range config.unmarshalFuncMap {
		bundle.RegisterUnmarshalFunc(format, unmarshalFunc)
	}

	for _, path := range config.translationFiles {
		_, err := bundle.LoadMessageFile(path)
		if err != nil {
			return err
		}
	}
	for _, translationFSFile := range config.translationFSFiles {
		for _, path := range translationFSFile.paths {
			_, err := bundle.LoadMessageFileFS(translationFSFile.fs, path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Get returns the translated message for the given message id.
//
// It uses the default language tag.
//
// Example:
//
//	message := i18n.Get("hello", i18n.Params{"name": "John"})
func Get(id string, opts ...any) string {
	return GetCtx(context.Background(), id, opts...)
}

// GetCtx returns the translated message for the given message id.
//
// It uses the language from the context. You can set the language to the context with i18n.Middleware.
// If the language is not found in the context, it uses the default language tag.
//
// Example:
//
//	message := i18n.GetCtx(ctx, "hello", i18n.Params{"name": "John"})
func GetCtx(ctx context.Context, id string, opts ...any) string {
	if bundle == nil {
		return "ERROR: i18n is not initialized"
	}

	cfg := newLocalizeConfig(opts...)
	localizeConfig := cfg.toI18nLocalizeConfig(id)

	var languages []string
	if cfg.language != "" {
		languages = append(languages, cfg.language)
	}
	lang := extractLanguageFunc(ctx)
	if lang != "" && !slices.Contains(languages, lang) {
		languages = append(languages, lang)
	}
	if !slices.Contains(languages, defaultLanguage.String()) {
		languages = append(languages, defaultLanguage.String())
	}

	localizer := i18n.NewLocalizer(bundle, languages...)
	message, err := localizer.Localize(localizeConfig)

	if message == "" {
		return missingTranslationHandler(id, err)
	}

	return message
}

// T is an alias for Get.
//
// Example:
//
//	message := i18n.T("hello", i18n.Params{"name": "John"})
func T(id string, opts ...any) string {
	return Get(id, opts...)
}

// TCtx is an alias for GetCtx.
//
// Example:
//
//	message := i18n.TCtx(ctx, "hello", i18n.Params{"name": "John"})
func TCtx(ctx context.Context, id string, opts ...any) string {
	return GetCtx(ctx, id, opts...)
}
