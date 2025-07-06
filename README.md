# fr0g.ai - AI-Powered Security Intelligence Platform

[![CI/CD](https://github.com/fr0g-vibe/fr0g-ai/workflows/CI/badge.svg)](https://github.com/fr0g-vibe/fr0g-ai/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/fr0g-vibe/fr0g-ai)](https://goreportcard.com/report/github.com/fr0g-vibe/fr0g-ai)

## Mission
Eliminate human-computer interaction vulnerabilities through AI-driven automated threat detection and response.

> *"You are the first layer of an artificial intelligence system designed to remove interactions between computers and humans."*

## Architecture

```
┌─────────────────┐    Events     ┌─────────────────┐    Analysis    ┌─────────────────┐
│   fr0g-ai-io    │◄─────────────►│ fr0g-ai-master- │◄──────────────►│   fr0g-ai-aip   │
│  (I/O Service)  │               │    control      │                │  (Core AI)      │
│   :8083/:9092   │               │ (Cognitive AI)  │                │   :8080/:9090   │
└─────────────────┘               └─────────────────┘                └─────────────────┘
         │                                 │                                   │
         ▼                                 ▼                                   ▼
   External I/O                    Intelligence Engine                   AI Processing
                                           │
                                           ▼
                                  ┌─────────────────┐
                                  │ fr0g-ai-bridge  │
                                  │  (Integration)  │
                                  │   :8082/:9091   │
                                  └─────────────────┘
                                           │
                                           ▼
                                     OpenWebUI API
```

### Components
- **fr0g-ai-registry**: Service discovery and health monitoring (port 8500) ⚠️ CRITICAL ISSUES
- **fr0g-ai-aip**: Core AI processing engine with file-based storage (ports 8080/9090) ⚠️ gRPC ISSUES
- **fr0g-ai-bridge**: Integration bridge connecting to external systems (ports 8082/9091) ⚠️ CRITICAL ISSUES
- **fr0g-ai-master-control**: Orchestration and cognitive processing engine (port 8081) ✅ OPERATIONAL
- **fr0g-ai-io**: Input/Output processing service for threat vector handling (ports 8083/9093) ⚠️ gRPC ISSUES
- **Communication**: High-performance gRPC inter-service communication ⚠️ NEEDS REPAIR
- **Storage**: Configurable storage backends (file system, Redis, future: database)

## Security Philosophy

### Baseline Principles
- **Trust no one, verify everything**
- **Zero human-computer interaction vulnerabilities**

### Threat Vectors Addressed
Human-computer interactions create vulnerabilities:
- **Email**: Unknown senders, malicious content, phishing attempts *(INTERCEPTED by ESMTP Processor)*
- **Phone**: Social engineering, voice spoofing, unauthorized access
- **Websites**: Drive-by downloads, malicious scripts, data harvesting
- **Software**: Flawed programs deployed by flawed humans creating systemic risks
- **Discord/Chat**: Social engineering through messaging platforms *(INTERCEPTED by Discord Processor)*

### Design Metaphor
*"A person walks into a tall building. They stop at the front desk and say 'I have a piece of mail that needs to get to a person.'"*

The fr0g.ai system acts as that intelligent front desk - intercepting, analyzing, and verifying all interactions before they reach their intended targets.

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/fr0g-vibe/fr0g-ai.git
cd fr0g-ai

# 2. Setup environment
cp .env.example .env
# Edit .env with your configuration

# 3. Initialize and build
make setup

# 4. Start services
docker-compose up -d

# 5. Verify health
make health
```

## Development

### Multi-Agent Development Environment

fr0g.ai uses a sophisticated multi-agent development system with tmux and aider for coordinated development across all components.

#### Quick Start
```bash
# Start the complete development environment
./start-fr0g-ai-dev.sh

# This creates a tmux session with 10 specialized agent windows:
# 0: Project-Lead    - Architecture & coordination
# 1: AIP            - Core AI service (fr0g-ai-aip)
# 2: Bridge         - Integration service (fr0g-ai-bridge)
# 3: MCP            - Cognitive engine (fr0g-ai-master-control)
# 4: IO             - I/O processing (fr0g-ai-io)
# 5: Config         - Configuration management
# 6: DevOps         - Infrastructure & deployment
# 7: Registry       - Service discovery (fr0g-ai-registry)
# 8: Build-Test     - Build automation
# 9: Shell          - Interactive shell
```

#### Agent Dispatch System

The Project Lead agent (window 0) can dispatch commands to other specialized agents:

```bash
# Command Format:
tmux send-keys -t fr0g-ai:WINDOW_NUMBER "COMMAND" C-m

# Example Dispatches:
tmux send-keys -t fr0g-ai:1 "Implement persona CRUD operations with validation" C-m
tmux send-keys -t fr0g-ai:2 "Add health check endpoint with metrics" C-m
tmux send-keys -t fr0g-ai:3 "Optimize learning rate for threat detection" C-m
tmux send-keys -t fr0g-ai:4 "Complete SMS processor integration" C-m
tmux send-keys -t fr0g-ai:7 "Enhance service discovery performance" C-m

# Shell Commands (windows 8-9):
tmux send-keys -t fr0g-ai:8 "make build-all" C-m
tmux send-keys -t fr0g-ai:9 "git status" C-m
```

#### Agent Specializations

Each agent has specific domain expertise and system prompts:

- **Project-Lead (0)**: Overall coordination, architecture decisions, cross-component integration
- **AIP Agent (1)**: Core AI processing, persona management, identity processing
- **Bridge Agent (2)**: External integrations, API gateway, protocol translation
- **MCP Agent (3)**: Cognitive intelligence, orchestration, conscious AI
- **IO Agent (4)**: Input/output processing, threat vector handling
- **Config Agent (5)**: Environment variables, shared config library, validation
- **DevOps Agent (6)**: Docker, deployment, CI/CD, infrastructure automation
- **Registry Agent (7)**: Service discovery, registration, health monitoring

#### Traditional Development (Alternative)

```bash
# Initialize development environment
make setup

# Build all services
make build-all

# Run services locally (separate terminals)
make run-aip     # Terminal 1: Core AI service
make run-bridge  # Terminal 2: Integration bridge

# Run tests
make test-all

# Code quality
make lint fmt
```

## Service Endpoints

### fr0g-ai-aip (Core AI)
- **HTTP**: http://localhost:8080
- **gRPC**: localhost:9090
- **Health**: http://localhost:8080/health

### fr0g-ai-registry (Service Discovery)
- **HTTP**: http://localhost:8500
- **Health**: http://localhost:8500/health
- **Service Registration**: http://localhost:8500/v1/agent/service/register
- **Service Discovery**: http://localhost:8500/v1/catalog/services

### fr0g-ai-bridge (Integration)
- **HTTP**: http://localhost:8082
- **gRPC**: localhost:9091
- **Health**: http://localhost:8082/health

### fr0g-ai-master-control (Cognitive Engine)
- **HTTP**: http://localhost:8081
- **Health**: http://localhost:8081/health

### fr0g-ai-io (Input/Output Processing)
- **HTTP**: http://localhost:8083
- **gRPC**: localhost:9092
- **Health**: http://localhost:8083/health

## Configuration

Key environment variables:
- `OPENWEBUI_API_KEY`: API key for OpenWebUI integration
- `FR0G_STORAGE_TYPE`: Storage backend type (file, database)
- `LOG_LEVEL`: Logging verbosity (debug, info, warn, error)

## Documentation

- [Architecture Overview](docs/architecture.md)
- [API Documentation](docs/api/)
- [Deployment Guide](docs/deployment.md)
- [Development Setup](docs/development.md)
- [AI Code Generation Guidelines](docs/ai-guidelines.md)

## AI Development Guidelines

### Essential Context Files for AI Sessions
When starting new AI coding sessions, always include these files:
- `README.md` (this file - project overview)
- `docker-compose.yml` (service configuration)
- `Makefile` (build commands)
- `.env.example` (configuration template)
- Component-specific TODO.md file for the service being worked on

### Component Boundaries
- **fr0g-ai-aip**: Core AI processing engine (ports 8080/9090)
- **fr0g-ai-bridge**: Integration bridge service (ports 8082/9091)  
- **fr0g-ai-master-control**: Orchestration and cognitive engine (port 8081)
- **fr0g-ai-io**: Input/Output processing service (ports 8083/9092)

### Repository Information
- **GitHub URL**: `https://github.com/fr0g-vibe/fr0g-ai`
- **Project Path**: `github.com/fr0g-vibe/fr0g-ai`
- **License**: GPL-3.0
- **Working Directory**: AI agents start in `/fr0g-ai` root (local clone of github.com/fr0g-vibe/fr0g-ai)
- **Module Navigation**: Must `cd` into component directory before Go commands
- **Subprojects**: All components exist under `github.com/fr0g-vibe/fr0g-ai/` path

### NO MOCKING POLICY
- **NEVER CREATE MOCKS**: Always implement real functionality, never create mock implementations
- **REPLACE EXISTING MOCKS**: If you find mock implementations, replace them with real working code
- **REAL INTEGRATIONS**: Always implement actual API calls, database connections, and service integrations
- **PRODUCTION READY**: All code must be production-ready, not placeholder or demo code

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Built by the fr0g.ai team*
