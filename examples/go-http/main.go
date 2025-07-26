package main

import (
	"fmt"
	"net/http"

	"github.com/afkdevs/go-i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const PORT = 3000

func main() {
	// Initialize the i18n package
	if err := i18n.Init(language.English,
		i18n.WithUnmarshalFunc("yaml", yaml.Unmarshal),
		i18n.WithTranslationFile("../../testdata/en.yaml", "../../testdata/id.yaml"),
	); err != nil {
		panic(err)
	}

	r := http.NewServeMux()
	// Register the handler and wrap it with i18n.Middleware
	r.Handle("/", i18n.Middleware(rootHandler()))
	r.Handle("/hello", i18n.Middleware(helloHandler()))

	fmt.Printf("Server running on http://localhost:%d \n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the translated message
		message := i18n.TCtx(r.Context(), "hello_world")

		// Write the message to the response
		w.Write([]byte(message))
	}
}

func helloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get name from query parameter
		name := r.URL.Query().Get("name")

		// Get the translated message with parameter
		message := i18n.TCtx(r.Context(), "hello", i18n.Params{"name": name})

		// Write the message to the response
		w.Write([]byte(message))
	}
}
