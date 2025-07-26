module github.com/afkdevs/go-i18n/examples/go-http

go 1.23.0

toolchain go1.24.5

require (
	github.com/afkdevs/go-i18n v0.0.0-0000000000000-000000000000
	golang.org/x/text v0.27.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/nicksnyder/go-i18n/v2 v2.6.0 // indirect

replace github.com/afkdevs/go-i18n => ../..
