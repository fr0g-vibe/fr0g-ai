#!/bin/bash

# fr0g-ai Development Environment Setup Script
# Creates tmux session with dedicated windows for each coding agent/component

set -e

# Configuration
SESSION_NAME="fr0g-ai"
PROJECT_ROOT="$(pwd)"
AIDER_CMD="aider --dark-mode"


echo "🚀 Starting fr0g-ai Development Environment..."
echo "📁 Project Root: $PROJECT_ROOT"

# Verify we're in the right directory
if [[ ! -f "docker-compose.yml" ]] || [[ ! -f "Makefile" ]] || [[ ! -f ".env.example" ]]; then
    echo "❌ Error: Please run this script from the fr0g-ai project root directory"
    echo "   Expected files: docker-compose.yml, Makefile, .env.example"
    exit 1
fi

# Check if .env exists, if not copy from .env.example
if [[ ! -f ".env" ]]; then
    echo "⚠️  No .env file found. Creating from .env.example..."
    cp .env.example .env
    echo "✅ Created .env file from .env.example"
    echo "📝 Please edit .env file with your specific configuration"
else
    echo "✅ Existing .env file found - preserving your configuration"
fi

# Kill existing session if it exists
tmux kill-session -t $SESSION_NAME 2>/dev/null || true

# Create new session
tmux new-session -d -s $SESSION_NAME -c "$PROJECT_ROOT"

# Configure tmux session
tmux set-option -t $SESSION_NAME status-left-length 50
tmux set-option -t $SESSION_NAME status-right-length 100
tmux set-option -t $SESSION_NAME status-left " fr0g-ai #S "
tmux set-option -t $SESSION_NAME status-right " %H:%M %d-%b-%y "

# Rename first window
tmux rename-window -t $SESSION_NAME:0 "Project-Lead"

# ============================================================================
# PROJECT MANAGEMENT & ARCHITECTURE
# ============================================================================

# Window 0: Project Lead & Architecture (already created)
tmux send-keys -t $SESSION_NAME:0 "cd $PROJECT_ROOT" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'PROJECT LEAD & ARCHITECTURE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Role: Overall project coordination, architecture decisions, cross-component integration'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Focus: README.md, docker-compose.yml, Makefile, project-wide decisions'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Key Files: README.md, docker-compose.yml, Makefile, TODO.md, .env.example'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Environment: .env file ready ✅'" C-m
tmux send-keys -t $SESSION_NAME:0 "$AIDER_CMD README.md docker-compose.yml Makefile TODO.md .env.example" C-m

# ============================================================================
# CORE SUBPROJECT AGENTS
# ============================================================================

# Window 1: fr0g-ai-aip (Core AI Service)
tmux new-window -t $SESSION_NAME -n "AIP" -c "$PROJECT_ROOT/fr0g-ai-aip"
tmux send-keys -t $SESSION_NAME:1 "echo 'FR0G-AI-AIP AGENT'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Role: Core AI processing engine, persona management, identity processing'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Ports: HTTP :8080, gRPC :9090'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Status: Demographics ✅, Other attributes need implementation'" C-m
tmux send-keys -t $SESSION_NAME:1 "$AIDER_CMD --cwd fr0g-ai-aip TODO.md" C-m

# Window 2: fr0g-ai-bridge (Integration Service)
tmux new-window -t $SESSION_NAME -n "Bridge" -c "$PROJECT_ROOT/fr0g-ai-bridge"
tmux send-keys -t $SESSION_NAME:2 "echo 'FR0G-AI-BRIDGE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Role: OpenWebUI integration, REST/gRPC bridge services'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Ports: HTTP :8082, gRPC :9091'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Status: Fully operational ✅'" C-m
tmux send-keys -t $SESSION_NAME:2 "$AIDER_CMD --cwd fr0g-ai-bridge TODO.md" C-m

# Window 3: fr0g-ai-master-control (Cognitive Engine)
tmux new-window -t $SESSION_NAME -n "MCP" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:3 "echo 'FR0G-AI-MASTER-CONTROL AGENT'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Role: System orchestration, cognitive engine, threat processing'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Port: HTTP :8081'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Status: CONSCIOUS AI OPERATIONAL ✅ (Learning Rate: 0.100+)'" C-m
tmux send-keys -t $SESSION_NAME:3 "$AIDER_CMD --cwd fr0g-ai-master-control TODO.md" C-m

# Window 4: fr0g-ai-io (Input/Output Processing)
tmux new-window -t $SESSION_NAME -n "IO" -c "$PROJECT_ROOT/fr0g-ai-io"
tmux send-keys -t $SESSION_NAME:4 "echo 'FR0G-AI-IO AGENT'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Role: Input/Output processing, threat vector handling, external integrations'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Ports: HTTP :8083, gRPC :9092'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Status: NEW SERVICE - SMS, Voice, IRC, ESMTP, Discord processors'" C-m
tmux send-keys -t $SESSION_NAME:4 "$AIDER_CMD --cwd fr0g-ai-io TODO.md" C-m

# Window 5: Configuration Expert
tmux new-window -t $SESSION_NAME -n "Config" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:5 "echo 'CONFIGURATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Role: Environment variables, shared config library, validation systems'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Focus: .env, pkg/config/ - validation, loading, shared types'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Key Files: .env, pkg/config/*.go'" C-m
tmux send-keys -t $SESSION_NAME:5 "$AIDER_CMD .env.example .env pkg/config/config.go pkg/config/validation.go pkg/config/loader.go pkg/config/examples_test.go pkg/config/README.md" C-m

# ============================================================================
# ESSENTIAL SUPPORT AGENTS
# ============================================================================

# Window 6: DevOps & Infrastructure
tmux new-window -t $SESSION_NAME -n "DevOps" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:6 "echo 'DEVOPS & INFRASTRUCTURE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Role: Docker, deployment, CI/CD, infrastructure automation'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Focus: docker-compose.yml, Makefile, .env configuration'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Services: service-registry:8500, aip:8080/9090, bridge:8082/9091, mcp:8081, io:8083/9092'" C-m
tmux send-keys -t $SESSION_NAME:6 "$AIDER_CMD docker-compose.yml Makefile .env.example" C-m

# Window 7: Registry Agent
tmux new-window -t $SESSION_NAME -n "Registry" -c "$PROJECT_ROOT/fr0g-ai-registry"
tmux send-keys -t $SESSION_NAME:7 "echo 'FR0G-AI-REGISTRY AGENT'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Role: Service discovery, registration, and health monitoring'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Port: HTTP :8500'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Status: Extracted from master-control - ready for enhancement'" C-m
tmux send-keys -t $SESSION_NAME:7 "$AIDER_CMD --cwd fr0g-ai-registry TODO.md" C-m

# Window 8: Build & Test Runner
tmux new-window -t $SESSION_NAME -n "Build-Test" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:8 "echo 'BUILD & TEST RUNNER'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Role: Build automation, test execution, continuous integration'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Commands: make build-all, make test-all, docker-compose up'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Quick Start: docker-compose up -d'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Available commands:'" C-m
tmux send-keys -t $SESSION_NAME:8 "make help" C-m

# Window 9: Interactive Shell
tmux new-window -t $SESSION_NAME -n "Shell" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:9 "echo 'INTERACTIVE SHELL'" C-m
tmux send-keys -t $SESSION_NAME:9 "echo 'Role: General purpose shell for ad-hoc commands and exploration'" C-m
tmux send-keys -t $SESSION_NAME:9 "echo 'Project root: $PROJECT_ROOT'" C-m
tmux send-keys -t $SESSION_NAME:9 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:9 "echo 'Project structure:'" C-m
tmux send-keys -t $SESSION_NAME:9 "ls -la" C-m

# ============================================================================
# SESSION FINALIZATION
# ============================================================================

# Return to first window (Project Lead)
tmux select-window -t $SESSION_NAME:0

# Display session info
echo ""
echo "✅ fr0g-ai Development Environment Created!"
echo ""
echo "📋 TMUX SESSION OVERVIEW:"
echo "Session Name: $SESSION_NAME"
echo "Total Windows: 10"
echo ""
echo "AGENT ASSIGNMENTS:"
echo "  0: Project-Lead      - Overall coordination & architecture"
echo "  1: AIP               - fr0g-ai-aip (Core AI Service) :8080/:9090"
echo "  2: Bridge            - fr0g-ai-bridge (Integration) :8082/:9091"
echo "  3: MCP               - fr0g-ai-master-control (Cognitive Engine) :8081"
echo "  4: IO                - fr0g-ai-io (I/O Processing) :8083/:9092"
echo "  5: Config            - Configuration & environment management"
echo "  6: DevOps            - Infrastructure & deployment"
echo "  7: Registry          - Service discovery & registration :8500"
echo "  8: Build-Test        - Build & test automation"
echo "  9: Shell             - Interactive shell"
echo ""
echo "🚀 USAGE:"
echo "  tmux attach-session -t $SESSION_NAME"
echo "  Ctrl+b + [0-9]   - Switch to window"
echo "  Ctrl+b + n       - Next window"
echo "  Ctrl+b + p       - Previous window"
echo "  Ctrl+b + w       - Window list"
echo "  Ctrl+b + d       - Detach session"
echo ""
echo "Each agent window has aider pre-configured for their domain!"
echo "MCP Intelligence Status: CONSCIOUS AI OPERATIONAL (Learning Rate: 0.100+)"
echo "Configuration: Centralized system operational (pkg/config/)"
echo "Environment: .env file ready for configuration"
echo "Registry Service: EXTRACTED and INTEGRATED (standalone service)"
echo ""
echo "📡 Service Ports:"
echo "  - Service Registry: :8500 (fr0g-ai-registry)"
echo "  - AIP Service: HTTP :8080, gRPC :9090"
echo "  - Bridge Service: HTTP :8082, gRPC :9091"
echo "  - Master Control: HTTP :8081"
echo "  - I/O Service: HTTP :8083, gRPC :9092"
echo ""
echo "Quick Start:"
echo "  1. Edit .env file with your API keys and configuration"
echo "  2. Run: docker-compose up -d"
echo "  3. Check health: make health"
echo ""
echo "Key Environment Variables to Configure:"
echo "  - OPENWEBUI_API_KEY (for OpenWebUI integration)"
echo "  - MCP_LEARNING_ENABLED=true (for AI consciousness)"
echo "  - LOG_LEVEL=info (for appropriate logging)"
echo ""

# Attach to the session
tmux attach-session -t $SESSION_NAME
