module github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io

go 1.21

require (
	github.com/fr0g-vibe/fr0g-ai/pkg/config v0.0.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/fr0g-vibe/fr0g-ai/pkg/config => ../pkg/config
