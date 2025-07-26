# fiber-i18n
This package provides a middleware for [Fiber](https://github.com/gofiber/fiber) that integrates with the go-i18n library for internationalization and localization.

## Installation
```bash
go get github.com/afkdevs/go-i18n/contrib/fiber-i18n
```

## Usage
```go
package main

import (
	"github.com/afkdevs/go-i18n/contrib/fiber-i18n"
	"github.com/gofiber/fiber/v2"
)

func main() {
    err := i18n.Init(language.English,
        i18n.WithTranslationFile("locales/en.yaml", "locales/id.yaml"),
    )
	if err != nil {
		log.Fatalf("failed to initialize i18n: %v", err)
	}
	app := fiber.New()
	app.Use(fiberi18n.New())

	app.Get("/", func(c *fiber.Ctx) error {
        // Using fiberi18n.TCtx to translate the text
        // It use the context from the Fiber Ctx
		return c.SendString(fiberi18n.TCtx(c, "test"))
	})
    app.Get("/hello/:name", func(c *fiber.Ctx) error {
        name := c.Params("name")
        // Using i18n.TCtx to translate the text with parameters
        // It use the user context from the Fiber Ctx
        return c.SendString(i18n.TCtx(c.UserContext(), "hello", i18n.Param("name", name)))
    })

	app.Listen(":3000")
}
```