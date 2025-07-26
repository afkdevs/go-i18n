package main

import (
	"fmt"
	"net/http"

	"github.com/afkdevs/go-i18n"
	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const PORT = 3000

func main() {
	// Initialize i18n
	if err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("../../testdata/en.yaml", "../../testdata/id.yaml"),
	); err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	// Register i18n middleware
	r.Use(i18n.NewMiddleware())

	// Register route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Get translation message with context
		message := i18n.TCtx(r.Context(), "hello_world")

		// Write response
		w.Write([]byte(message))
	})
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		// Get name from query parameter
		name := chi.URLParam(r, "name")

		// Get translation message with context
		message := i18n.TCtx(r.Context(), "hello_name", i18n.Params{"name": name})

		// Write response
		w.Write([]byte(message))
	})

	fmt.Printf("Server running on http://localhost:%d \n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}
