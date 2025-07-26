# go-i18n

[![Go](https://github.com/afkdevs/go-i18n/actions/workflows/ci.yml/badge.svg)](https://github.com/afkdevs/go-i18n/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/afkdevs/go-i18n)](https://goreportcard.com/report/github.com/afkdevs/go-i18n)
[![codecov](https://codecov.io/gh/afkdevs/go-i18n/graph/badge.svg?token=DPEMJ3DgRX)](https://codecov.io/gh/afkdevs/go-i18n)
[![GoDoc](https://pkg.go.dev/badge/github.com/afkdevs/go-i18n)](https://pkg.go.dev/github.com/afkdevs/go-i18n)
[![Go Version](https://img.shields.io/github/go-mod/go-version/afkdevs/go-i18n)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

go-i18n is a simple internationalization library for Go.
This wraps [nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n) and adds useful features for easier integration and customization.


## Installation
```bash
go get -u github.com/afkdevs/go-i18n
```

## Features

- [x] Support for translation files in **YAML**, **JSON**, and **TOML** formats
- [x] Simple string translation
- [x] Context-based translation
- [x] Parameterized translation
- [x] Fallback for missing translations
- [x] Customizable language extraction from context

## Usage

### Create Translation File
**`locales/en.yaml`**
```yaml
hello: Hello
hello_name: Hello, {{.name}}
hello_name_age: Hello, {{.name}}. You are {{.age}} years old
```

**`locales/id.yaml`**
```yaml
hello: Halo
hello_name: Halo, {{.name}}
hello_name_age: Halo, {{.name}}. Kamu berumur {{.age}} tahun
```

### Initialize i18n
```go
package main

import (
	"log"

	"github.com/afkdevs/go-i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func main() {
	err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("locales/en.yaml", "locales/id.yaml"),
	)
	if err != nil {
		log.Fatalf("failed to initialize i18n: %v", err)
	}
}
```

### Translate your text

#### Simple translation

```go
msg := i18n.T("hello")

// Or with a specific language
msg = i18n.T("hello", i18n.Lang("id"))
```

#### Translation with parameters

```go
// Single parameter
msg := i18n.T("hello_name", i18n.Param("name", "John"))

// Multiple parameters
msg := i18n.T("hello_name_age", i18n.Params{
	"name": "John",
	"age":  20,
})

// Or using a map
msg := i18n.T("hello_name_age", map[string]any{
	"name": "John",
	"age":  20,
})
```

## Context Translation

Use `TCtx` to translate using a `context.Context`, which is helpful for request-scoped translations.

Example with **chi**:

```go
package main

import (
	"net/http"

	"github.com/afkdevs/go-i18n"
	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func main() {
	if err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("locales/en.yaml", "locales/id.yaml"),
	); err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	// Automatically inject language context
	r.Use(i18n.NewMiddleware())

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		message := i18n.TCtx(r.Context(), "hello")
		_, _ = w.Write([]byte(message))
	})
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		message := i18n.TCtx(r.Context(), "hello_name", i18n.Params{"name": name})
		_, _ = w.Write([]byte(message))
	})

	http.ListenAndServe(":3000", r)
}
```

## Fallback for Missing Translations

You can set a global configuration for missing translations by using `i18n.WithMissingTranslationHandler` when initializing i18n.

```go
package main

import (
	"fmt"
	"log"

	"github.com/afkdevs/go-i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func main() {
	if err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("locales/en.yaml", "locales/id.yaml"),
		i18n.WithMissingTranslationHandler(missingTranslationHandler),
	); err != nil {
		panic(err)
	}
}

func missingTranslationHandler(id string, _ error) string {
	return fmt.Sprintf("ERROR: missing translation for %q", id)
}
```

## Contributing

Contributions are welcome!  
If you’d like to make major changes, please open an issue first to discuss your ideas.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details