module github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io

go 1.21

require (
	github.com/fr0g-vibe/fr0g-ai/pkg/config v0.0.0
	google.golang.org/grpc v1.58.0
	google.golang.org/protobuf v1.31.0
)

replace github.com/fr0g-vibe/fr0g-ai/pkg/config => ../pkg/config
