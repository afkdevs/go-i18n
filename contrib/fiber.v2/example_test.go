package fiber_test

import (
	"log"

	"github.com/afkdevs/go-i18n"
	fiberi18n "github.com/afkdevs/go-i18n/contrib/fiber.v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func Example() {
	if err := i18n.Init(language.English,
		i18n.WithTranslationFile("../../testdata/en.yaml", "../../testdata/id.yaml"),
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
	); err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(fiberi18n.New())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(fiberi18n.TCtx(c, "test"))
	})
	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Query("name")
		return c.SendString(fiberi18n.GetCtx(c, "hello_name", i18n.Param("name", name)))
	})

	log.Fatal(app.Listen(":3000"))
}

func Example_customConfig() {
	if err := i18n.Init(language.English,
		i18n.WithTranslationFile("../../testdata/en.yaml", "../../testdata/id.yaml"),
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
	); err != nil {
		panic(err)
	}

	app := fiber.New()

	// Custom configuration for i18n middleware
	// Here we set the language based on a query parameter
	// You can also use headers or any other method to determine the language
	app.Use(fiberi18n.New(
		fiberi18n.WithLanguageHandler(func(c *fiber.Ctx) string {
			return c.Query("lang", "en")
		}),
	))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(fiberi18n.TCtx(c, "test"))
	})
	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Query("name")
		return c.SendString(fiberi18n.GetCtx(c, "hello_name", i18n.Param("name", name)))
	})

	log.Fatal(app.Listen(":3000"))
}
