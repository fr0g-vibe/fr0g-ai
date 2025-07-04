# fr0g.ai - AI-Powered Security Intelligence Platform

[![CI/CD](https://github.com/fr0g-vibe/fr0g-ai/workflows/CI/badge.svg)](https://github.com/fr0g-vibe/fr0g-ai/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/fr0g-vibe/fr0g-ai)](https://goreportcard.com/report/github.com/fr0g-vibe/fr0g-ai)

## 🎯 Mission
Eliminate human-computer interaction vulnerabilities through AI-driven automated threat detection and response.

> *"You are the first layer of an artificial intelligence system designed to remove interactions between computers and humans."*

## 🏗️ Architecture

```
┌─────────────────┐    gRPC     ┌─────────────────┐
│   fr0g-ai-aip   │◄───────────►│ fr0g-ai-bridge  │
│   (Core AI)     │             │  (Integration)  │
│   :8080/:9090   │             │   :8081/:9091   │
└─────────────────┘             └─────────────────┘
         │                               │
         ▼                               ▼
   File Storage                    OpenWebUI API
```

### Components
- **fr0g-ai-aip**: Core AI processing engine with file-based storage
- **fr0g-ai-bridge**: Integration bridge connecting to external systems (OpenWebUI)
- **fr0g-ai-master-control**: Orchestration and cognitive processing engine
- **ESMTP Threat Vector Interceptor**: Email intelligence gathering and threat analysis
- **Communication**: High-performance gRPC inter-service communication
- **Storage**: Configurable storage backends (file system, future: database)

## 🚨 Security Philosophy

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

## 🚀 Quick Start

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

## 🛠️ Development

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

## 📊 Service Endpoints

### fr0g-ai-aip (Core AI)
- **HTTP**: http://localhost:8080
- **gRPC**: localhost:9090
- **Health**: http://localhost:8080/health

### fr0g-ai-bridge (Integration)
- **HTTP**: http://localhost:8081
- **gRPC**: localhost:9091
- **Health**: http://localhost:8081/health

## 🔧 Configuration

Key environment variables:
- `OPENWEBUI_API_KEY`: API key for OpenWebUI integration
- `FR0G_STORAGE_TYPE`: Storage backend type (file, database)
- `LOG_LEVEL`: Logging verbosity (debug, info, warn, error)

## 📚 Documentation

- [Architecture Overview](docs/architecture.md)
- [API Documentation](docs/api/)
- [Deployment Guide](docs/deployment.md)
- [Development Setup](docs/development.md)
- [AI Code Generation Guidelines](docs/ai-guidelines.md)

## 🤖 AI Development Guidelines

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

### Repository Information
- **GitHub URL**: `https://github.com/fr0g-vibe/fr0g-ai`
- **Project Path**: `github.com/fr0g-vibe/fr0g-ai`
- **License**: GPL-3.0
- **Working Directory**: AI agents start in `/fr0g-ai` root (local clone of github.com/fr0g-vibe/fr0g-ai)
- **Module Navigation**: Must `cd` into component directory before Go commands
- **Subprojects**: All components exist under `github.com/fr0g-vibe/fr0g-ai/` path

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Built with 🐸 by the fr0g.ai team*
