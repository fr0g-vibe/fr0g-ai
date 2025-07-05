module github.com/fr0g-vibe/fr0g-ai

go 1.22

replace github.com/fr0g-vibe/fr0g-ai-aip => ./fr0g-ai-aip

replace github.com/fr0g-vibe/fr0g-ai-bridge => ./fr0g-ai-bridge
module github.com/fr0g-vibe/fr0g-ai

go 1.21

require (
	github.com/gorilla/mux v1.8.0
	google.golang.org/grpc v1.58.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
)
