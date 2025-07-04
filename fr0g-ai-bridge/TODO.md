# fr0g-ai-bridge TODO

## 🤖 AI CODE GENERATION GUIDELINES - BRIDGE COMPONENT

### 📋 ESSENTIAL CONTEXT FILES FOR THIS COMPONENT
**ALWAYS ADD THESE FILES TO AI CHAT CONTEXT:**
- `README.md` (project overview and component boundaries)
- `docker-compose.yml` (service configuration)
- `Makefile` (build commands)
- `.env.example` (configuration template)
- `fr0g-ai-bridge/TODO.md` (THIS FILE - current status)
- `fr0g-ai-bridge/internal/api/validation.go` (current implementation)

### 🚨 COMPONENT BOUNDARY RULES
- **FOCUS AREA**: Only work on `fr0g-ai-bridge/` directory and files
- **SERVICE ROLE**: Integration bridge between OpenWebUI and fr0g-ai-aip
- **PORTS**: HTTP :8082, gRPC :9091 (configured in docker-compose)
- **DEPENDENCIES**: Calls fr0g-ai-aip via gRPC, integrates with OpenWebUI

### ⚠️ CROSS-COMPONENT INTERACTION RULES
- **DO NOT** edit files in `fr0g-ai-aip/` or `fr0g-ai-master-control/` directories
- **DO NOT** modify other components' TODO.md files
- **ASK FIRST** if you need fr0g-ai-aip protobuf definitions or gRPC interfaces
- **ASK FIRST** if you need to modify shared files (docker-compose.yml, Makefile, etc.)
- **BE AWARE** that you consume fr0g-ai-aip services but don't implement them

### 🏗️ PROJECT STRUCTURE RULES
- **Repository URL**: Always use `https://github.com/fr0g-vibe/fr0g-ai`
- **Project Path**: All Go imports use `github.com/fr0g-vibe/fr0g-ai/` prefix
- **Working Directory**: AI agents start in `/fr0g-ai` root directory (local clone)
- **Module Navigation**: MUST `cd fr0g-ai-bridge` before running Go commands
- **Service Ports**: HTTP :8082, gRPC :9091 (configured in docker-compose)
- **Subproject Path**: This component exists at `github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge`

### 🚫 PROTOBUF GENERATION RULES
- **NEVER EDIT**: Do not manually edit any `.pb.go` files marked "DO NOT EDIT"
- **Use Build Commands**: Always use `make proto` or `protoc` commands for protobuf generation
- **Generated Files**: Treat all `.pb.go` files as build artifacts, not source code
- **Import Generated**: Import generated protobuf code, never attempt to create it manually
- **AIP Protobuf**: Use existing protobuf definitions from fr0g-ai-aip, do not recreate them

### 🌉 BRIDGE SERVICE SPECIFIC GUIDELINES
- **Primary Role**: Integration bridge between OpenWebUI and fr0g-ai-aip
- **Communication**: REST API inbound, gRPC outbound to AIP service
- **Validation**: Comprehensive request/response validation required
- **Error Handling**: Graceful degradation when AIP service unavailable

### 🔌 INTEGRATION PATTERNS
- **OpenWebUI Client**: Implement HTTP client with retry logic
- **AIP gRPC Client**: Use connection pooling and health checking
- **Service Discovery**: Register with service registry on startup
- **Configuration**: Environment-based configuration with validation

### 📡 API DESIGN STANDARDS
- **REST Endpoints**: Follow OpenAPI 3.0 specification
- **gRPC Services**: Use protobuf definitions from AIP service
- **Request Validation**: Validate all inputs before processing
- **Response Format**: Consistent JSON response structure
- **Error Codes**: Use appropriate HTTP status codes

### 🛡️ SECURITY IMPLEMENTATION
- **API Authentication**: Implement API key validation
- **Request Sanitization**: Clean all user inputs
- **Rate Limiting**: Implement per-client rate limiting
- **CORS**: Configure CORS for web client access

## High Priority - Core Functionality

### Validation System
- [ ] **CRITICAL**: Implement comprehensive request validation beyond role checking (only isValidRole() exists)
- [ ] **CRITICAL**: Add persona context validation (referenced in docker-compose but missing)
- [ ] **CRITICAL**: Implement message content validation and sanitization
- [ ] **CRITICAL**: Add model parameter validation (temperature, max_tokens, etc.)

### Missing Core Handlers
- [ ] **CRITICAL**: Implement chat completion handlers (completely missing - no handlers directory exists)
- [ ] **CRITICAL**: Implement gRPC service handlers (GRPC_PORT configured but no gRPC code)
- [ ] **CRITICAL**: Create actual REST API handlers in internal/api/ (only validation.go exists)
- [ ] **CRITICAL**: Implement OpenWebUI client integration (OPENWEBUI_BASE_URL configured but no client)

### Service Integration Missing
- [ ] **CRITICAL**: AIP_GRPC_ENDPOINT configured but no gRPC client implementation
- [ ] **CRITICAL**: Service registry integration missing (REGISTRY_URL configured but not used)
- [ ] **CRITICAL**: Health check endpoint missing (referenced in docker-compose)

### Chat Completion Service
- [ ] Implement actual chat completion handlers (currently stubbed)
- [ ] Add streaming response support
- [ ] Implement conversation context management
- [ ] Add model selection and routing logic

### OpenWebUI Integration
- [ ] Complete OpenWebUI client implementation
- [ ] Add authentication handling
- [ ] Implement error handling and retries
- [ ] Add connection pooling and timeout management

### gRPC Service Implementation
- [ ] Implement all gRPC service methods
- [ ] Add proper error handling and status codes
- [ ] Implement streaming gRPC endpoints
- [ ] Add gRPC middleware for logging/auth

### Service Discovery Integration
- [ ] Implement service registry client
- [ ] Add automatic service registration/deregistration
- [ ] Implement health checks with dependency status
- [ ] Add service discovery for AIP service connection

## Medium Priority - Features

### Request/Response Management
- [ ] Implement request validation middleware
- [ ] Add response caching layer
- [ ] Implement request/response logging
- [ ] Add request tracing and correlation IDs

### Security & Authentication
- [ ] Implement API key authentication
- [ ] Add CORS configuration management
- [ ] Implement rate limiting per client
- [ ] Add request sanitization and validation

### Persona Integration
- [ ] Implement persona-aware chat completions
- [ ] Add persona context injection
- [ ] Implement persona switching mid-conversation
- [ ] Add persona performance tracking

### Error Handling & Resilience
- [ ] Implement circuit breaker pattern
- [ ] Add retry logic with exponential backoff
- [ ] Implement graceful degradation
- [ ] Add comprehensive error logging

## Low Priority - Nice to Have

### Monitoring & Observability
- [ ] Add metrics collection (Prometheus)
- [ ] Implement distributed tracing
- [ ] Add performance monitoring
- [ ] Create health check dashboard

### Developer Experience
- [ ] Add OpenAPI/Swagger documentation
- [ ] Implement comprehensive test suite
- [ ] Add mock servers for development
- [ ] Create integration test framework

### Advanced Features
- [ ] Implement conversation history storage
- [ ] Add conversation analytics
- [ ] Implement A/B testing framework
- [ ] Add conversation export/import

## Technical Debt

### Code Organization
- [ ] Refactor main.go - extract server setup logic
- [ ] Implement proper dependency injection
- [ ] Add comprehensive error types
- [ ] Improve configuration management

### Testing
- [ ] Add unit tests for all handlers
- [ ] Implement integration tests
- [ ] Add load testing framework
- [ ] Create end-to-end test suite

### Documentation
- [ ] Add API documentation
- [ ] Create integration guides
- [ ] Write troubleshooting documentation
- [ ] Add performance tuning guides

## Immediate Actions Needed

### Missing Implementations
- [ ] **CRITICAL**: Implement actual REST handlers in internal/handlers/
- [ ] **CRITICAL**: Implement gRPC service handlers in internal/grpc/
- [ ] **CRITICAL**: Create OpenWebUI client implementation in internal/clients/
- [ ] **HIGH**: Add proper validation beyond role checking

### New Framework Implementation Tasks
- [x] **COMPLETED**: Create handlers framework directory
- [x] **COMPLETED**: Create clients framework directory  
- [x] **COMPLETED**: Create grpc framework directory
- [ ] **URGENT**: Implement handlers/chat completion endpoints
- [ ] **URGENT**: Implement clients/openwebui integration
- [ ] **URGENT**: Implement grpc/bridge service methods
- [ ] **HIGH**: Implement clients/aip gRPC client
