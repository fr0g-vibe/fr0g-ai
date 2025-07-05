module github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip

go 1.21

require (
	github.com/fr0g-vibe/fr0g-ai/pkg/config v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/fr0g-vibe/fr0g-ai/pkg/config => ../pkg/config
