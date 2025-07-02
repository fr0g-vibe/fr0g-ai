# Development Guide

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Git with submodule support
- Make
- Protocol Buffers compiler (protoc)
- golangci-lint (optional, for code quality)

## Getting Started

### 1. Clone and Setup

```bash
# Clone with submodules
git clone --recursive https://github.com/fr0g-vibe/fr0g-ai.git
cd fr0g-ai

# Setup development environment
make setup
```

### 2. Configure Environment

```bash
# Copy and edit environment file
cp .env.example .env
# Edit .env with your specific configuration
```

### 3. Development Workflow

#### Option A: Docker Development (Recommended)
```bash
# Start all services with Docker
make dev

# View logs
make docker-logs

# Stop services
docker-compose down
```

#### Option B: Local Development
```bash
# Terminal 1: Start AIP service
make run-aip

# Terminal 2: Start Bridge service  
make run-bridge

# Terminal 3: Check health
make health
```

## Project Structure

```
fr0g-ai/
├── .github/workflows/     # CI/CD pipelines
├── docs/                  # Documentation
├── fr0g-ai-aip/          # Core AI service (submodule)
├── fr0g-ai-bridge/       # Integration bridge (submodule)
├── data/                  # Local data storage
├── config/                # Configuration files
├── logs/                  # Application logs
├── docker-compose.yml     # Container orchestration
├── Makefile              # Build automation
└── .env.example          # Environment template
```

## Development Commands

### Building
```bash
make build-all          # Build all services
make docker-build-all   # Build Docker images
```

### Testing
```bash
make test-all          # Run all tests
make lint              # Run linters
make fmt               # Format code
```

### Running Services
```bash
make run-aip           # Run AIP service locally
make run-bridge        # Run Bridge service locally
make health            # Check service health
```

### Docker Operations
```bash
make dev               # Start development environment
make docker-logs       # View container logs
make docker-clean      # Clean containers and volumes
```

## Service Development

### fr0g-ai-aip (Core AI Service)

**Key Files:**
- `main.go` - Service entry point
- `internal/` - Internal packages
- `proto/` - gRPC definitions
- `Dockerfile` - Container configuration

**Development:**
```bash
cd fr0g-ai-aip
make build-with-grpc    # Build with gRPC support
make test               # Run tests
make proto              # Generate protobuf files
```

### fr0g-ai-bridge (Integration Bridge)

**Key Files:**
- `main.go` - Service entry point
- `internal/` - Internal packages
- `proto/` - gRPC definitions
- `Dockerfile` - Container configuration

**Development:**
```bash
cd fr0g-ai-bridge
make build-with-grpc    # Build with gRPC support
make test               # Run tests
make proto              # Generate protobuf files
```

## API Development

### gRPC Services

Both services expose gRPC APIs for inter-service communication:

- **AIP Service**: `localhost:9090`
- **Bridge Service**: `localhost:9091`

### HTTP APIs

Both services also expose HTTP APIs:

- **AIP Service**: `http://localhost:8080`
- **Bridge Service**: `http://localhost:8081`

### Health Checks

All services implement health check endpoints:

```bash
curl http://localhost:8080/health  # AIP service
curl http://localhost:8081/health  # Bridge service
```

## Debugging

### Logs

```bash
# View all container logs
make docker-logs

# View specific service logs
docker-compose logs -f fr0g-ai-aip
docker-compose logs -f fr0g-ai-bridge
```

### Local Debugging

```bash
# Run with debug logging
LOG_LEVEL=debug make run-aip
LOG_LEVEL=debug make run-bridge
```

### gRPC Debugging

Use tools like `grpcurl` to test gRPC endpoints:

```bash
# List services
grpcurl -plaintext localhost:9090 list

# Call methods
grpcurl -plaintext -d '{"message": "test"}' localhost:9090 service.Method
```

## Testing

### Unit Tests
```bash
make test-all           # All tests
cd fr0g-ai-aip && make test    # AIP tests only
cd fr0g-ai-bridge && make test # Bridge tests only
```

### Integration Tests
```bash
# Start services and run integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### Manual Testing
```bash
# Health checks
make health

# API testing with curl
curl -X POST http://localhost:8080/api/endpoint -d '{"data": "test"}'
```

## Code Quality

### Linting
```bash
make lint              # Run all linters
make fmt               # Format all code
```

### Pre-commit Hooks
```bash
# Install pre-commit hooks (if using pre-commit)
pre-commit install
```

## Troubleshooting

### Common Issues

1. **Submodule Issues**
   ```bash
   make update-submodules  # Force update submodules
   ```

2. **Build Failures**
   ```bash
   make clean             # Clean build artifacts
   make deps              # Reinstall dependencies
   ```

3. **Docker Issues**
   ```bash
   make docker-clean      # Clean Docker resources
   docker system prune -f # Clean Docker system
   ```

4. **Port Conflicts**
   - Check if ports 8080, 8081, 9090, 9091 are available
   - Modify port mappings in docker-compose.yml if needed

### Getting Help

1. Check the [Architecture Documentation](architecture.md)
2. Review service-specific README files in submodules
3. Check GitHub Issues for known problems
4. Review Docker logs for error messages

## Contributing

### Development Workflow

1. Create a feature branch
2. Make changes in appropriate submodule
3. Test changes locally
4. Update documentation if needed
5. Submit pull request

### Code Standards

- Follow Go best practices
- Use gofmt for formatting
- Write tests for new functionality
- Update documentation for API changes
- Follow semantic versioning for releases

### Commit Messages

Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation
- `refactor:` for code refactoring
- `test:` for test additions
