# i18n
i18n is a simple internationalization library for Golang.
It provides a way to translate strings into multiple languages.
This library is wrapper for [go-i18n](https://github.com/nicksnyder/go-i18n) with some additional features.

## Installation
```bash
go get -u github.com/ahmadfaizk/i18n
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
# locale/en.yaml
hello: Hello
hello_name: Hello, {{.name}}
hello_name_age: Hello, {{.name}}. You are {{.age}} years old
```

```yaml
# locale/id.yaml
hello: Halo
hello_name: Halo, {{.name}}
hello_name_age: Halo, {{.name}}. Kamu berumur {{.age}} tahun
```

### Initialize i18n
```go
import (
	"github.com/ahmadfaizk/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

i18n.Init(language.English,
    i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
    i18n.WithTranslationFile("locale/en.yaml", "locale/id.yaml"),
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
fmt.Println(i18n.T("hello_name_age", i18n.Params(i18n.M{"name": "John", "age": 20})))
// Hello, John. You are 20 years old
fmt.Println(i18n.T("hello_name_age", i18n.Lang("id"), i18n.Params(i18n.M{"name": "John", "age": 20})))
// Halo, John. Kamu berumur 20 tahun
```

## Examples with Chi

```go
package main

import (
	"github.com/ahmadfaizk/i18n"
	"github.com/go-chi/chi"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"net/http"
)

func main() {
	if err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("locale/en.yaml", "locale/id.yaml"),
	); err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(i18n.Middleware)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		message := i18n.TCtx(ctx, "hello")
		_, _ = w.Write([]byte(message))
	})
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := chi.URLParam(r, "name")
		message := i18n.TCtx(ctx, "hello_name", i18n.Param("name", name))
		_, _ = w.Write([]byte(message))
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details