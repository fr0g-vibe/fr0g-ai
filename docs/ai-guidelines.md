# fr0g.ai AI Code Generation Guidelines

## Project Overview
- **Repository**: https://github.com/fr0g-vibe/fr0g-ai
- **Project Path**: `github.com/fr0g-vibe/fr0g-ai`
- **License**: GPL-3.0
- **Architecture**: Microservices with gRPC communication
- **Language**: Go 1.21+
- **Subproject Structure**: All components exist under `github.com/fr0g-vibe/fr0g-ai/` path

## 🏗️ Component Architecture

### Service Boundaries
```
┌─────────────────┐    gRPC     ┌─────────────────┐
│   fr0g-ai-aip   │◄───────────►│ fr0g-ai-bridge  │
│   (Core AI)     │             │  (Integration)  │
│   :8080/:9090   │             │   :8081/:9091   │
└─────────────────┘             └─────────────────┘
         │                               │
         ▼                               ▼
   File Storage                    OpenWebUI API
         ▲
         │
┌─────────────────┐
│fr0g-ai-master-  │
│    control      │
│     :8081       │
└─────────────────┘
```

### Component Responsibilities
- **fr0g-ai-aip**: Core AI processing, persona/identity management, data storage
- **fr0g-ai-bridge**: Integration layer, OpenWebUI communication, request routing
- **fr0g-ai-master-control**: Orchestration, cognitive processing, threat analysis

## 🚨 AI Agent Boundaries

### Component Isolation Rules
1. **Work only within assigned component directory**
2. **Never edit other components' files without explicit permission**
3. **Ask before modifying shared files** (docker-compose.yml, Makefile, etc.)
4. **Request information about other components rather than assuming**

### 🚫 NO MOCKING POLICY - CRITICAL RULE
1. **NEVER CREATE MOCK IMPLEMENTATIONS**: Always implement real, working functionality
2. **REPLACE EXISTING MOCKS**: If you encounter mock code, replace it with real implementations
3. **REAL API INTEGRATIONS**: Always implement actual HTTP clients, gRPC clients, database connections
4. **PRODUCTION-READY CODE**: All implementations must be production-ready, not placeholders
5. **NO STUB METHODS**: Never create empty methods or functions that return fake data
6. **REAL ERROR HANDLING**: Implement proper error handling for actual failure scenarios
7. **ACTUAL EXTERNAL CALLS**: Make real calls to external services, databases, APIs

### Cross-Component Communication
- Use gRPC for service-to-service communication
- Respect defined interfaces and contracts
- Ask for API documentation when integrating with other services

## SETUP️ Development Standards

### Go Development Rules
- **Working Directory**: Always start in `/fr0g-ai` root (local clone of github.com/fr0g-vibe/fr0g-ai)
- **Project Path**: All Go imports use `github.com/fr0g-vibe/fr0g-ai/` prefix
- **Module Navigation**: `cd` into component before Go commands
- **Go Version**: Use Go 1.21+ features
- **Error Handling**: Never ignore errors with `_`
- **Context**: Pass `context.Context` as first parameter when needed

### Protobuf Generation Rules
- **NEVER EDIT**: Do not manually edit any `.pb.go` files that contain "DO NOT EDIT" comments
- **Use Build Commands**: Always use `make proto` or component-specific protobuf build commands
- **Generated Files**: Treat all `.pb.go` files as generated artifacts, not source code
- **Protobuf Sources**: Only edit `.proto` files, never the generated Go code
- **Build First**: Run protobuf generation before building Go code that depends on it

### Implementation Standards - NO MOCKING
- **REAL IMPLEMENTATIONS ONLY**: Never create mock, stub, or placeholder implementations
- **ACTUAL FUNCTIONALITY**: Every function must perform its intended real-world operation
- **PRODUCTION QUALITY**: All code must be ready for production deployment
- **REAL EXTERNAL INTEGRATIONS**: Implement actual API calls, not simulated responses
- **PROPER ERROR HANDLING**: Handle real error conditions, not just return nil
- **COMPLETE FEATURES**: Implement full functionality, not partial or demo versions

### Code Style
- Use `gofmt` and `goimports`
- Follow golangci-lint rules
- Use descriptive names, avoid abbreviations
- Implement comprehensive error handling

### Testing Requirements
- Unit tests for all business logic
- Integration tests for external dependencies
- Mock external services in tests
- Follow table-driven test patterns

## 📁 Directory Structure Standards

```
{component}/
├── cmd/                    # Main applications
├── internal/               # Private application code
├── pkg/                   # Public library code (if any)
├── configs/               # Configuration files
├── scripts/               # Build scripts
├── tests/                 # Integration tests
└── docs/                  # Component documentation
```

## Build and Deployment

### Make Targets
- Use Makefile for common operations
- Component-specific targets should be in component directories
- Root Makefile coordinates multi-component operations

### Docker
- Each component has its own Dockerfile
- Multi-stage builds for optimization
- Health checks required for all services

### Configuration
- Use `.env` files for local development
- Environment variables for production
- Validate configuration on startup

## Security Standards

### Input Validation
- Validate all external inputs
- Sanitize user-provided data
- Use typed validation where possible

### Error Handling
- Don't leak sensitive information in errors
- Log security events appropriately
- Implement proper error boundaries

### Dependencies
- Keep dependencies updated
- Use go mod for dependency management
- Regular security audits

## 📝 Documentation Requirements

### Code Documentation
- Godoc format for public APIs
- Inline comments for complex logic
- README.md for each component

### API Documentation
- OpenAPI specs for REST APIs
- Protobuf documentation for gRPC
- Integration examples

## 🔄 Git Workflow

### Commit Standards
- Use conventional commit format
- Include component prefix in commits
- Reference issues in commit messages

### Branch Strategy
- Feature branches for development
- Component-specific branch naming
- Clean commit history

## Deployment Considerations

### Service Discovery
- Register with service registry
- Implement health checks
- Handle service dependencies

### Monitoring
- Structured logging with correlation IDs
- Metrics for key operations
- Distributed tracing support

### Scaling
- Stateless design where possible
- Resource limits and requests
- Graceful shutdown handling
