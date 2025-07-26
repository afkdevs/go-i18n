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
	// Register the handler
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/hello", helloHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: i18n.NewMiddleware()(r), // Use i18n middleware
	}

	fmt.Printf("Server running on http://localhost:%d \n", PORT)
	server.ListenAndServe()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Get the translated message
	message := i18n.TCtx(r.Context(), "hello_world")

	// Write the message to the response
	w.Write([]byte(message))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Get name from query parameter
	name := r.URL.Query().Get("name")

	// Get the translated message with parameter
	message := i18n.TCtx(r.Context(), "hello_name", i18n.Params{"name": name})

	// Write the message to the response
	w.Write([]byte(message))
}
