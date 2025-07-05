# fr0g-ai-registry TODO

## AI CODE GENERATION GUIDELINES - REGISTRY COMPONENT

### ESSENTIAL CONTEXT FILES FOR THIS COMPONENT
**ALWAYS ADD THESE FILES TO AI CHAT CONTEXT:**
- `README.md` (project overview and component boundaries)
- `docker-compose.yml` (service configuration)
- `Makefile` (build commands)
- `.env.example` (configuration template)
- `fr0g-ai-registry/TODO.md` (THIS FILE - current status)

### COMPONENT BOUNDARY RULES
- **FOCUS AREA**: Only work on `fr0g-ai-registry/` directory and files
- **SERVICE ROLE**: Service discovery, registration, and health monitoring
- **PORTS**: HTTP :8500 (configured in docker-compose)
- **DEPENDENCIES**: Standalone service, provides discovery to all other services

### CROSS-COMPONENT INTERACTION RULES
- **DO NOT** edit files in other component directories
- **DO NOT** modify other components' TODO.md files
- **ASK FIRST** if you need to modify shared files (docker-compose.yml, Makefile, etc.)
- **BE AWARE** that all other services depend on this registry for service discovery

### PROJECT STRUCTURE RULES
- **Repository URL**: Always use `https://github.com/fr0g-vibe/fr0g-ai`
- **Project Path**: All Go imports use `github.com/fr0g-vibe/fr0g-ai/` prefix
- **Working Directory**: AI agents start in `/fr0g-ai` root directory (local clone)
- **Module Navigation**: MUST `cd fr0g-ai-registry` before running Go commands
- **Service Ports**: HTTP :8500 (configured in docker-compose)
- **Subproject Path**: This component exists at `github.com/fr0g-vibe/fr0g-ai/fr0g-ai-registry`

### NO MOCKING POLICY - REGISTRY COMPONENT
- **NEVER CREATE MOCKS**: Always implement real functionality, never mock implementations
- **REAL SERVICE DISCOVERY**: Implement actual service registration and discovery
- **REAL HEALTH CHECKING**: Implement actual health monitoring, not fake status
- **REAL PERSISTENCE**: Implement actual service state persistence
- **PRODUCTION READY**: All registry functionality must be production-ready

### CODE QUALITY REQUIREMENTS - REGISTRY COMPONENT
- **MANDATORY LINTING**: Always run `make lint` before committing any code changes
- **ZERO LINT ERRORS**: All code must pass golangci-lint without errors or warnings
- **FIX BEFORE COMMIT**: Never commit code that fails linting - fix all issues first
- **LINT EARLY**: Run `make lint` frequently during development, not just at the end
- **SHARED CONFIG**: Use centralized configuration from `pkg/config/` to avoid import errors

### SEARCH/REPLACE BLOCK RULES - REGISTRY COMPONENT
- **QUADRUPLE BACKTICKS**: Always use ```` as fences, never triple backticks ```
- **FULL FILE PATH**: Use complete file path alone on first line, no formatting
- **EXACT MATCHING**: SEARCH section must match existing content character-for-character
- **CONCISE BLOCKS**: Keep blocks small, include only changing lines plus context
- **UNIQUE MATCHING**: Include enough surrounding lines for unique identification
- **MULTIPLE BLOCKS**: Use separate blocks for multiple changes in same file

### 🚨 CRITICAL SAFETY RULES - REGISTRY COMPONENT 🚨
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

### CENTRALIZED CONFIGURATION RULES - REGISTRY COMPONENT
- **MANDATORY**: Use `pkg/config/` for ALL configuration and validation needs
- **NO LOCAL CONFIG**: Never create registry-specific config/validation libraries
- **EXTEND SHARED**: Embed shared config types, add registry-specific fields as needed
- **VALIDATION STANDARD**: Use `sharedconfig.ValidationErrors` for all validation
- **IMPORT PATTERN**: Always `import sharedconfig "pkg/config"`
- **CONTRIBUTE SPECIALIZED**: Add registry-specific validation to `pkg/config/` when needed
- **NO DUPLICATION**: Never reimplement port, timeout, or other validation already in shared config
- **LOADER USAGE**: Use `sharedconfig.NewLoader()` for configuration loading

## STATUS: EXTRACTION COMPLETE - SERVICE FULLY INTEGRATED

### ✅ COMPLETED: Service Extraction and Integration
- [x] **EXTRACTED**: Registry server implementation from master-control
- [x] **EXTRACTED**: Main entry point and HTTP server setup
- [x] **EXTRACTED**: Consul-compatible API endpoints
- [x] **EXTRACTED**: Service registration/deregistration functionality
- [x] **EXTRACTED**: Health checking and monitoring
- [x] **EXTRACTED**: Docker containerization support
- [x] **COMPLETED**: Go module setup with proper dependencies
- [x] **COMPLETED**: Makefile integration with build system
- [x] **COMPLETED**: Clean compilation and binary generation
- [x] **COMPLETED**: Integration with main project build system
- [x] **COMPLETED**: Tmux development environment integration
- [x] **COMPLETED**: Service startup script integration
- [x] **COMPLETED**: Service shutdown script integration
- [x] **COMPLETED**: Health monitoring integration

### ✅ COMPLETED: Core Functionality
- [x] **OPERATIONAL**: HTTP API server on port 8500
- [x] **OPERATIONAL**: Service registration endpoint (/v1/agent/service/register)
- [x] **OPERATIONAL**: Service deregistration endpoint (/v1/agent/service/deregister/{id})
- [x] **OPERATIONAL**: Service discovery endpoint (/v1/catalog/services)
- [x] **OPERATIONAL**: Health check endpoint (/health)
- [x] **OPERATIONAL**: Service health monitoring (/v1/health/service/{id})
- [x] **OPERATIONAL**: Build system integration (make build-all includes registry)
- [x] **OPERATIONAL**: Standalone service binary (bin/fr0g-ai-registry)
- [x] **OPERATIONAL**: Development environment (tmux window 8)
- [x] **OPERATIONAL**: Service lifecycle management (start/stop scripts)

## High Priority - Service Enhancement

### Service Discovery Features
- [ ] **HIGH**: Add service load balancing and routing
- [ ] **HIGH**: Implement service health check automation
- [ ] **HIGH**: Add service metadata and tagging support
- [ ] **MEDIUM**: Implement service versioning support
- [ ] **MEDIUM**: Add service dependency tracking

### Persistence and Reliability
- [ ] **HIGH**: Add persistent storage for service registry data
- [ ] **HIGH**: Implement service state recovery on restart
- [ ] **MEDIUM**: Add backup and restore functionality
- [ ] **LOW**: Implement distributed registry clustering

### Monitoring and Observability
- [ ] **HIGH**: Add comprehensive metrics collection
- [ ] **HIGH**: Implement service discovery analytics
- [ ] **MEDIUM**: Add alerting for service failures
- [ ] **LOW**: Create registry dashboard

## Medium Priority - Advanced Features

### API Enhancements
- [ ] **MEDIUM**: Add pagination for service listings
- [ ] **MEDIUM**: Implement advanced service filtering
- [ ] **MEDIUM**: Add bulk service operations
- [ ] **LOW**: Implement GraphQL API support

### Security and Authentication
- [ ] **MEDIUM**: Add API authentication and authorization
- [ ] **MEDIUM**: Implement service-to-service authentication
- [ ] **LOW**: Add audit logging for all operations
- [ ] **LOW**: Implement rate limiting per client

### Integration Features
- [ ] **MEDIUM**: Add webhook notifications for service changes
- [ ] **MEDIUM**: Implement service mesh integration
- [ ] **LOW**: Add Kubernetes service discovery integration
- [ ] **LOW**: Create CLI tools for registry management

## Low Priority - Nice to Have

### Developer Experience
- [ ] **LOW**: Add OpenAPI/Swagger documentation
- [ ] **LOW**: Implement comprehensive test suite
- [ ] **LOW**: Create SDK for common languages
- [ ] **LOW**: Add development seed data

### Performance Optimization
- [ ] **LOW**: Implement caching for frequently accessed data
- [ ] **LOW**: Add connection pooling
- [ ] **LOW**: Optimize service lookup algorithms
- [ ] **LOW**: Add performance profiling

## Technical Debt

### Code Organization
- [ ] **MEDIUM**: Implement proper dependency injection
- [ ] **MEDIUM**: Add comprehensive error handling
- [ ] **LOW**: Refactor large functions and improve modularity
- [ ] **LOW**: Add comprehensive documentation

### Testing
- [ ] **HIGH**: Add unit tests for all registry functions
- [ ] **HIGH**: Implement integration tests with other services
- [ ] **MEDIUM**: Add load testing framework
- [ ] **LOW**: Create end-to-end test suite

### Configuration
- [ ] **MEDIUM**: Implement configuration validation
- [ ] **MEDIUM**: Add environment-specific configs
- [ ] **LOW**: Implement configuration hot-reloading
- [ ] **LOW**: Add feature flags

## IMMEDIATE ACTIONS - PHASE 1: SERVICE EXTRACTION

### ✅ COMPLETED: Framework Extraction
1. ✅ **EXTRACTED**: Service registry server from master-control
2. ✅ **EXTRACTED**: HTTP API endpoints and routing
3. ✅ **EXTRACTED**: Service registration and discovery logic
4. ✅ **EXTRACTED**: Health monitoring functionality
5. ✅ **EXTRACTED**: Docker containerization support

### 🎯 IMMEDIATE NEXT PRIORITIES: Service Testing and Enhancement
1. **HIGH**: Test registry service with all fr0g.ai services (integration testing)
2. **HIGH**: Verify service registration/discovery workflow with real services
3. **HIGH**: Fix service configuration issues (port conflicts, storage validation)
4. **MEDIUM**: Add persistent storage for service data
5. **MEDIUM**: Implement automated health checking
6. **MEDIUM**: Add comprehensive metrics and monitoring
7. **LOW**: Enhance service discovery with load balancing
8. **LOW**: Add authentication and security features

## SUCCESS METRICS

### Service Metrics
- **Uptime**: 99.9% availability target
- **Response Time**: <50ms for service lookups
- **Throughput**: Handle 1000+ service operations per minute
- **Discovery Accuracy**: 100% accurate service discovery

### Integration Metrics
- **Service Registration**: All fr0g.ai services auto-register
- **Health Monitoring**: Real-time health status for all services
- **Discovery Performance**: Sub-millisecond service lookups
- **Reliability**: Zero data loss during service restarts

### Current Implementation Status
- **HTTP API**: ✅ OPERATIONAL (Consul-compatible endpoints)
- **Service Registration**: ✅ OPERATIONAL (register/deregister)
- **Service Discovery**: ✅ OPERATIONAL (catalog and health APIs)
- **Health Monitoring**: ✅ OPERATIONAL (automated health checks)
- **Docker Integration**: ✅ OPERATIONAL (containerized deployment)
- **Build System**: ✅ OPERATIONAL (Makefile and go.mod complete, builds successfully)
- **Binary Generation**: ✅ OPERATIONAL (bin/fr0g-ai-registry executable)
- **Project Integration**: ✅ OPERATIONAL (integrated with make build-all)
- **Development Environment**: ✅ OPERATIONAL (tmux window 8 configured)
- **Service Lifecycle**: ✅ OPERATIONAL (start/stop scripts integrated)
- **Health Monitoring**: ✅ OPERATIONAL (make health includes registry)
- **Persistence**: ⏳ PENDING (currently in-memory only)
- **Metrics**: ⏳ PENDING (basic health only)
- **Service Testing**: ⏳ PENDING (needs integration testing with other services)

### Quality Metrics
- **Test Coverage**: >80% code coverage target
- **Documentation**: Complete API and integration documentation
- **Security**: All API endpoints secured and validated
- **Monitoring**: Comprehensive metrics and alerting in place
