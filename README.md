# go-i18n

[![Go](https://github.com/afkdevs/go-i18n/actions/workflows/ci.yml/badge.svg)](https://github.com/afkdevs/go-i18n/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/afkdevs/go-i18n)](https://goreportcard.com/report/github.com/afkdevs/go-i18n)
[![GoDoc](https://pkg.go.dev/badge/github.com/afkdevs/go-i18n)](https://pkg.go.dev/github.com/afkdevs/go-i18n)
[![Go Version](https://img.shields.io/github/go-mod/go-version/afkdevs/go-i18n)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

go-i18n is a simple internationalization library for Golang.
It provides a way to translate strings into multiple languages.
This library is wrapper for [go-i18n](https://github.com/nicksnyder/go-i18n) with some additional features.

## Table of Contents
- [Installation](#installation)
- [Features](#features)
- [Usage](#usage)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Installation
```bash
go get -u github.com/afkdevs/go-i18n
```

## Features
- [x] Support translation file in YAML, JSON, and TOML format
- [x] Translation
- [x] Context translation
- [x] Parametrized translation
- [x] Missing translation Fallback
- [x] Custom extract language from Context

## Usage

### Create Translation File
```yaml
# locales/en.yaml
hello: Hello
hello_name: Hello, {{.name}}
hello_name_age: Hello, {{.name}}. You are {{.age}} years old
```

```yaml
# locales/id.yaml
hello: Halo
hello_name: Halo, {{.name}}
hello_name_age: Halo, {{.name}}. Kamu berumur {{.age}} tahun
```

### Initialize i18n
```go
import (
	"github.com/afkdevs/go-i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

i18n.Init(language.English,
    i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
    i18n.WithTranslationFile("locales/en.yaml", "locales/id.yaml"),
)
```

### Translate your text
```go
fmt.Println(i18n.T("hello"))
// Hello
fmt.Println(i18n.T("hello", i18n.Lang("id")))
// Halo
fmt.Println(i18n.T("hello_name", i18n.Param("name", "John")))
// Hello, John
fmt.Println(i18n.T("hello_name", i18n.Lang("id"), i18n.Param("name", "John")))
// Halo, John
fmt.Println(i18n.T("hello_name_age", i18n.Params{"name": "John", "age": 20}))
// Hello, John. You are 20 years old
fmt.Println(i18n.T("hello_name_age", i18n.Lang("id"), i18n.Params{"name": "John", "age": 20}))
// Halo, John. Kamu berumur 20 tahun
```

## Examples
See [examples/](/examples/) for a variety of examples.
```go
package main

import (
	"github.com/afkdevs/go-i18n"
	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"net/http"
)

func main() {
	if err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("locales/en.yaml", "locales/id.yaml"),
	); err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(i18n.Middleware)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		message := i18n.TCtx(r.Context(), "hello")
		_, _ = w.Write([]byte(message))
	})
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		message := i18n.TCtx(r.Context(), "hello_name", i18n.Params{"name": name})
		_, _ = w.Write([]byte(message))
	})

	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}
```
```bash
curl http://localhost:3000/hello # Hello
curl http://localhost:3000/hello/John # Hello, John
curl -H "Accept-Language: id" http://localhost:3000/hello # Halo
curl -H "Accept-Language: id" http://localhost:3000/hello/John # Halo, John
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details