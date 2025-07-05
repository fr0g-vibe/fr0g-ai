module github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control

go 1.21

require (
	github.com/fr0g-vibe/fr0g-ai/pkg/config v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
	gopkg.in/yaml.v2 v2.4.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

replace github.com/fr0g-vibe/fr0g-ai/pkg/config => ../pkg/config
