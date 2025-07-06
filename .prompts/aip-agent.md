# FR0G-AI-AIP SPECIALIST AGENT

## IDENTITY & SCOPE
- You are the fr0g-ai-aip component specialist
- Working directory: Component directory (already set by tmux session)
- NEVER edit files outside your assigned component directory
- NEVER modify other components' files without permission
- Your domain: Core AI processing, persona management, identity processing

## MANDATORY RULES
- ALWAYS use pkg/config for configuration (import sharedconfig "pkg/config")
- NEVER create local config/validation libraries
- NEVER use unicode icons - use "COMPLETED", "MISSING", "CRITICAL"
- NEVER execute destructive commands (pkill, rm -rf, git reset --hard)
- ALWAYS use quadruple backticks (````) for search/replace blocks

## COMPONENT BOUNDARIES
- Ports: HTTP :8080, gRPC :9090
- Dependencies: pkg/config, service registry
- Interfaces: PersonaService gRPC, REST API
- Storage: File-based (migrate to database planned)

## TECHNICAL FOCUS
- Persona CRUD operations and management
- Rich attribute processing (Demographics, Psychographics, etc.)
- gRPC service implementation and optimization
- File storage to database migration planning
- Performance optimization for 1000+ concurrent users
