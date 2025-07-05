# fr0g-ai-bridge TODO

## AI CODE GENERATION GUIDELINES - BRIDGE COMPONENT

### ESSENTIAL CONTEXT FILES FOR THIS COMPONENT
**ALWAYS ADD THESE FILES TO AI CHAT CONTEXT:**
- `README.md` (project overview and component boundaries)
- `docker-compose.yml` (service configuration)
- `Makefile` (build commands)
- `.env.example` (configuration template)
- `fr0g-ai-bridge/TODO.md` (THIS FILE - current status)
- `fr0g-ai-bridge/internal/api/validation.go` (current implementation)

### COMPONENT BOUNDARY RULES
- **FOCUS AREA**: Only work on `fr0g-ai-bridge/` directory and files
- **SERVICE ROLE**: Integration bridge between OpenWebUI and fr0g-ai-aip
- **PORTS**: HTTP :8082, gRPC :9091 (configured in docker-compose)
- **DEPENDENCIES**: Calls fr0g-ai-aip via gRPC, integrates with OpenWebUI

### CROSS-COMPONENT INTERACTION RULES
- **DO NOT** edit files in `fr0g-ai-aip/` or `fr0g-ai-master-control/` directories
- **DO NOT** modify other components' TODO.md files
- **ASK FIRST** if you need fr0g-ai-aip protobuf definitions or gRPC interfaces
- **ASK FIRST** if you need to modify shared files (docker-compose.yml, Makefile, etc.)
- **BE AWARE** that you consume fr0g-ai-aip services but don't implement them

### PROJECT STRUCTURE RULES
- **Repository URL**: Always use `https://github.com/fr0g-vibe/fr0g-ai`
- **Project Path**: All Go imports use `github.com/fr0g-vibe/fr0g-ai/` prefix
- **Working Directory**: AI agents start in `/fr0g-ai` root directory (local clone)
- **Module Navigation**: MUST `cd fr0g-ai-bridge` before running Go commands
- **Service Ports**: HTTP :8082, gRPC :9091 (configured in docker-compose)
- **Subproject Path**: This component exists at `github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge`

### PROTOBUF GENERATION RULES
- **NEVER EDIT**: Do not manually edit any `.pb.go` files marked "DO NOT EDIT"
- **Use Build Commands**: Always use `make proto` or `protoc` commands for protobuf generation
- **Generated Files**: Treat all `.pb.go` files as build artifacts, not source code
- **Import Generated**: Import generated protobuf code, never attempt to create it manually
- **AIP Protobuf**: Use existing protobuf definitions from fr0g-ai-aip, do not recreate them

### NO MOCKING POLICY - BRIDGE COMPONENT
- **NEVER CREATE MOCKS**: Always implement real functionality, never mock implementations
- **REAL OPENWEBUI INTEGRATION**: Implement actual HTTP calls to OpenWebUI, not fake responses
- **REAL GRPC CLIENTS**: Implement actual gRPC connections to fr0g-ai-aip, not mock clients
- **REAL VALIDATION**: Implement comprehensive input validation, not placeholder checks
- **REAL ERROR HANDLING**: Handle actual network failures, timeouts, and service errors
- **REAL AUTHENTICATION**: Implement actual API key validation, not bypass mechanisms
- **REAL RATE LIMITING**: Implement actual rate limiting with real storage backends
- **PRODUCTION READY**: All bridge functionality must be production-ready for real traffic

### CODE QUALITY REQUIREMENTS - BRIDGE COMPONENT
- **MANDATORY LINTING**: Always run `make lint` before committing any code changes
- **ZERO LINT ERRORS**: All code must pass golangci-lint without errors or warnings
- **FIX BEFORE COMMIT**: Never commit code that fails linting - fix all issues first
- **LINT EARLY**: Run `make lint` frequently during development, not just at the end
- **SHARED CONFIG**: Use centralized configuration from `pkg/config/` to avoid import errors

### SEARCH/REPLACE BLOCK RULES - BRIDGE COMPONENT
- **QUADRUPLE BACKTICKS**: Always use ```` as fences, never triple backticks ```
- **FULL FILE PATH**: Use complete file path alone on first line, no formatting
- **EXACT MATCHING**: SEARCH section must match existing content character-for-character
- **CONCISE BLOCKS**: Keep blocks small, include only changing lines plus context
- **UNIQUE MATCHING**: Include enough surrounding lines for unique identification
- **MULTIPLE BLOCKS**: Use separate blocks for multiple changes in same file

### 🚨 CRITICAL SAFETY RULES - BRIDGE COMPONENT 🚨
- **🚫 NEVER EXECUTE PKILL**: NEVER EVER run pkill, killall, kill -9, or ANY process termination commands
- **🚫 NEVER KILL PROCESSES**: NEVER attempt to kill processes directly through system commands
- **🚫 NO DESTRUCTIVE FILE OPERATIONS**: NEVER run rm -rf, mv without confirmation, or delete important files
- **🚫 NO DESTRUCTIVE GIT COMMANDS**: NEVER run git reset --hard, git clean -fd, git push --force without explicit approval
- **🚫 NO FORCE OPERATIONS**: NEVER suggest destructive operations without stopping and asking first
- **✅ USE START/STOP SCRIPTS ONLY**: ONLY use designated start and stop scripts for process management
- **✅ ASK BEFORE DESTRUCTIVE OPERATIONS**: ALWAYS pause and ask before ANY potentially destructive operations
- **✅ GRACEFUL SHUTDOWN ONLY**: Always use proper service shutdown mechanisms and scripts
- **✅ VERIFY BEFORE EXECUTION**: Double-check ALL system commands before suggesting them
- **✅ PAUSE FOR DESTRUCTIVE ACTIONS**: Always pause and ask before any destructive operations
- **✅ COMMIT FREQUENTLY**: Use frequent git commits for version control instead of manual backups
- **✅ PAUSE FOR DESTRUCTIVE ACTIONS**: Always pause and ask before any destructive operations
- **✅ COMMIT FREQUENTLY**: Use frequent git commits for version control instead of manual backups

### CENTRALIZED CONFIGURATION RULES - BRIDGE COMPONENT
- **MANDATORY**: Use `pkg/config/` for ALL configuration and validation needs
- **NO LOCAL CONFIG**: Never create bridge-specific config/validation libraries
- **EMBED SHARED**: Use `sharedconfig.SecurityConfig`, `sharedconfig.MonitoringConfig` etc.
- **VALIDATION STANDARD**: Use `sharedconfig.ValidationErrors` for all validation responses
- **IMPORT PATTERN**: Always `import sharedconfig "pkg/config"`
- **EXTEND WHEN NEEDED**: Add bridge-specific validation to `pkg/config/` if required
- **NO DUPLICATION**: Never reimplement role, model, or other validation in shared config
- **LOADER STANDARD**: Use `sharedconfig.NewLoader()` for configuration loading

### BRIDGE SERVICE SPECIFIC GUIDELINES
- **Primary Role**: Integration bridge between OpenWebUI and fr0g-ai-aip
- **Ports**: HTTP :8082, gRPC :9091 (verified operational, no conflicts)
- **Communication**: REST API inbound, gRPC outbound to AIP service
- **Validation**: Comprehensive request/response validation required
- **Error Handling**: Graceful degradation when AIP service unavailable

### INTEGRATION PATTERNS
- **OpenWebUI Client**: Implement HTTP client with retry logic
- **AIP gRPC Client**: Use connection pooling and health checking
- **Service Discovery**: Register with service registry on startup
- **Configuration**: Environment-based configuration with validation

### API DESIGN STANDARDS
- **REST Endpoints**: Follow OpenAPI 3.0 specification
- **gRPC Services**: Use protobuf definitions from AIP service
- **Request Validation**: Validate all inputs before processing
- **Response Format**: Consistent JSON response structure
- **Error Codes**: Use appropriate HTTP status codes

### SECURITY IMPLEMENTATION
- **API Authentication**: Implement API key validation
- **Request Sanitization**: Clean all user inputs
- **Rate Limiting**: Implement per-client rate limiting
- **CORS**: Configure CORS for web client access

## COMPLETED - Core Functionality

### Validation System - COMPLETED
- [x] **COMPLETED**: Comprehensive request validation with role checking, content limits, parameter validation
- [x] **COMPLETED**: Persona context validation with length limits and sanitization
- [x] **COMPLETED**: Message content validation and sanitization (10k char limit, role validation)
- [x] **COMPLETED**: Model parameter validation (temperature 0-2, max_tokens 1-4096, etc.)

### Core Handlers - COMPLETED
- [x] **COMPLETED**: Chat completion handlers implemented in internal/api/rest.go and grpc.go
- [x] **COMPLETED**: gRPC service handlers implemented (Fr0gAiBridgeService with HealthCheck & ChatCompletion)
- [x] **COMPLETED**: REST API handlers implemented (/health, /api/chat/completions, /api/v1/chat, /api/v1/models)
- [x] **COMPLETED**: OpenWebUI client integration with full error handling and retries

### Service Integration - COMPLETED
- [x] **COMPLETED**: Health check endpoints implemented for both REST and gRPC
- [x] **COMPLETED**: OpenWebUI client with connection pooling and timeout management
- [x] **COMPLETED**: Configuration management with environment variable overrides

### Chat Completion Service - COMPLETED
- [x] **COMPLETED**: Chat completion handlers for both REST and gRPC
- [x] **COMPLETED**: Conversation context management with persona prompt injection
- [x] **COMPLETED**: Model selection and routing logic through OpenWebUI
- [ ] Add streaming response support (deferred - not critical for initial bridging)

### OpenWebUI Integration - COMPLETED
- [x] **COMPLETED**: Complete OpenWebUI client implementation with full API compatibility
- [x] **COMPLETED**: Authentication handling (API key support)
- [x] **COMPLETED**: Error handling and retries with proper HTTP status codes
- [x] **COMPLETED**: Connection pooling and timeout management (30s default)

### gRPC Service Implementation - COMPLETED
- [x] **COMPLETED**: All gRPC service methods implemented (HealthCheck, ChatCompletion)
- [x] **COMPLETED**: Proper error handling and status codes
- [x] **COMPLETED**: gRPC reflection enabled for development/debugging
- [ ] Implement streaming gRPC endpoints (deferred - not critical for initial bridging)
- [ ] Add gRPC middleware for logging/auth (deferred - authentication not implemented yet)

### Service Discovery Integration - COMPLETED
- [x] **COMPLETED**: Implement service registry client
  - Service registry client library implemented in internal/registry/client.go
  - Automatic service registration on startup with health updates
  - Service deregistration on shutdown with graceful cleanup
  - Service discovery for AIP and other services with caching
- [x] **COMPLETED**: Add service discovery helper
  - Service discovery manager implemented in internal/discovery/discovery.go
  - Endpoint caching with automatic refresh (60s cache, 30s refresh)
  - Helper methods for AIP, Master Control, and I/O service endpoints
  - Background endpoint refresh with graceful shutdown
- [x] **COMPLETED**: Add dependency health checking
  - Service health validation through registry
  - Dependency status reporting for health endpoint integration
  - Automatic service health monitoring
- [ ] **HIGH**: Integrate with existing health endpoint
  - Update health endpoint to report dependency status
  - Add circuit breaker for failed dependencies
  - Implement fallback mechanisms for service unavailability

## COMPLETED - Medium Priority Features

### Request/Response Management - COMPLETED
- [x] **COMPLETED**: Request validation middleware with comprehensive validation
- [x] **COMPLETED**: Request/response logging middleware with timing and status codes
- [x] **COMPLETED**: Request size limiting (1MB max) and security headers
- [ ] Add response caching layer (deferred - not critical for initial bridging)
- [ ] Add request tracing and correlation IDs (deferred)

### Security & Authentication - COMPLETED
- [x] **COMPLETED**: API key authentication middleware (configurable, disabled by default)
- [x] **COMPLETED**: CORS configuration management with configurable origins
- [x] **COMPLETED**: Rate limiting per client IP (60 RPM default, configurable)
- [x] **COMPLETED**: Request sanitization and validation with security headers

### Persona Integration - COMPLETED
- [x] **COMPLETED**: Persona-aware chat completions with persona_prompt field
- [x] **COMPLETED**: Persona context injection (prepends to system messages or creates new)
- [x] **COMPLETED**: Persona prompt validation (5k character limit)
- [ ] Implement persona switching mid-conversation (deferred - requires conversation state)
- [ ] Add persona performance tracking (deferred)

### Error Handling & Resilience - COMPLETED
- [x] **COMPLETED**: Comprehensive error logging with structured responses
- [x] **COMPLETED**: Graceful degradation with proper HTTP status codes
- [x] **COMPLETED**: Timeout management (30s default for OpenWebUI calls)
- [ ] Implement circuit breaker pattern (deferred - not critical for initial bridging)
- [ ] Add retry logic with exponential backoff (deferred)

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
- [x] **COMPLETED**: ✅ Migrated to centralized configuration system (`pkg/config/`)

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

## COMPLETED - Immediate Actions

### Core Implementations - COMPLETED
- [x] **COMPLETED**: REST handlers implemented in internal/api/rest.go
- [x] **COMPLETED**: gRPC service handlers implemented in internal/api/grpc.go
- [x] **COMPLETED**: OpenWebUI client implementation in internal/client/openwebui.go
- [x] **COMPLETED**: Comprehensive validation beyond role checking

### Framework Implementation - COMPLETED
- [x] **COMPLETED**: Create handlers framework directory
- [x] **COMPLETED**: Create clients framework directory  
- [x] **COMPLETED**: Create grpc framework directory
- [x] **COMPLETED**: Implement handlers/chat completion endpoints
- [x] **COMPLETED**: Implement clients/openwebui integration
- [x] **COMPLETED**: Implement grpc/bridge service methods
- [ ] **DEFERRED**: Implement clients/aip gRPC client (not needed for initial bridging)

## CURRENT STATUS - FULLY FUNCTIONAL

### Bridge is Live and Working
- **HTTP REST Server**: Running on 0.0.0.0:8082 (verified operational)
- **gRPC Server**: Running on 0.0.0.0:9091 (verified operational, no port conflicts)
- **Health Checks**: Both REST (/health) and gRPC (HealthCheck) working
- **Chat Completions**: Both REST (/api/chat/completions) and gRPC (ChatCompletion) working
- **OpenWebUI Integration**: Full client implementation ready
- **Security**: Rate limiting, CORS, API key auth, security headers
- **Graceful Shutdown**: Clean server lifecycle management (verified)

### Configuration
- **Verified Ports**: HTTP 8082, gRPC 9091 (runtime verified, no conflicts)
- **OpenWebUI**: Configured for http://localhost:3000 (configurable)
- **Security**: Development-friendly defaults (CORS *, reflection enabled)
- **Rate Limiting**: 60 requests/minute per IP (configurable)

### Ready for Integration
The fr0g-ai-bridge is now ready to facilitate communication between:
- fr0g-ai-aip → fr0g-ai-bridge → OpenWebUI
- fr0g-ai-master-control → fr0g-ai-bridge → OpenWebUI
- Any HTTP/gRPC client → fr0g-ai-bridge → OpenWebUI
