module github.com/fr0g-vibe/fr0g-ai

go 1.23.0

toolchain go1.24.3

require github.com/gorilla/mux v1.8.1

// Exclude problematic internal imports
exclude github.com/fr0g-vibe/fr0g-ai/internal/grpc v0.0.0
