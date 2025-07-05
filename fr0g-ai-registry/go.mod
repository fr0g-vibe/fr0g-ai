module github.com/fr0g-vibe/fr0g-ai/fr0g-ai-registry

go 1.23.0

toolchain go1.24.3

require (
	github.com/fr0g-vibe/fr0g-ai/pkg/config v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
)

require (
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/fr0g-vibe/fr0g-ai/pkg/config => ../pkg/config
