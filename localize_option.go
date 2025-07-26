package i18n

import (
	"reflect"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Params is an alias for map[string]any. It is used to set template data for the message.
//
// Example:
//
//	i18n.T("hello", i18n.Params{"name": "John", "age": 30})
type Params map[string]any

type localizeConfig struct {
	params         map[string]any
	defaultMessage string
	language       string
}

func newLocalizeConfig(opts ...any) *localizeConfig {
	c := &localizeConfig{
		params: make(map[string]any),
	}
	for _, opt := range opts {
		reflectValue := reflect.ValueOf(opt)
		if reflectValue.Kind() == reflect.Map {
			for _, key := range reflectValue.MapKeys() {
				c.params[key.String()] = reflectValue.MapIndex(key).Interface()
			}
		} else if reflectValue.Kind() == reflect.Func {
			// Check if the function is a LocalizeOption
			if localizeOpt, ok := opt.(LocalizeOption); ok {
				localizeOpt(c)
				continue
			}
		}
	}
	return c
}

func (c *localizeConfig) toI18nLocalizeConfig(id string) *i18n.LocalizeConfig {
	localizeConfig := &i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: c.params,
	}
	if c.defaultMessage != "" {
		localizeConfig.DefaultMessage = &i18n.Message{
			ID:    id,
			Other: c.defaultMessage,
		}
	}
	return localizeConfig
}

// LocalizeOption is a function that configures the localizeConfig.
type LocalizeOption func(*localizeConfig)

// Param set single value of template data for the message.
//
// Example:
//
//	i18n.T("hello", i18n.Param("name", "John"))
func Param(key string, value any) LocalizeOption {
	return func(c *localizeConfig) {
		c.params[key] = value
	}
}

// Lang sets the language for the message.
//
// Example:
//
//	i18n.T("hello", i18n.Lang("id"))
func Lang(language string) LocalizeOption {
	return func(c *localizeConfig) {
		c.language = language
	}
}

// Default sets the default message for the message.
//
// It is used when the message is not found.
//
// Example:
//
//	i18n.T("hello", i18n.Default("Hello, {{.name}}!"), i18n.Param("name", "John")))
func Default(defaultMessage string) LocalizeOption {
	return func(c *localizeConfig) {
		c.defaultMessage = defaultMessage
	}
}
