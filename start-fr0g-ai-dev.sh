#!/bin/bash

# fr0g-ai Development Environment Setup Script
# Creates tmux session with dedicated windows for each coding agent/component

set -e

# Configuration
SESSION_NAME="fr0g-ai"
PROJECT_ROOT="$(pwd)"
AIDER_CMD="aider --no-auto-commits --dark-mode"

# Colors for tmux status
export TMUX_STATUS_BG="colour234"
export TMUX_STATUS_FG="colour39"

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
tmux set-option -t $SESSION_NAME status-bg $TMUX_STATUS_BG
tmux set-option -t $SESSION_NAME status-fg $TMUX_STATUS_FG
tmux set-option -t $SESSION_NAME status-left-length 50
tmux set-option -t $SESSION_NAME status-right-length 100
tmux set-option -t $SESSION_NAME status-left "#[fg=colour39,bg=colour234] fr0g-ai #[fg=colour234,bg=colour39] #S "
tmux set-option -t $SESSION_NAME status-right "#[fg=colour39] %H:%M %d-%b-%y "

# Rename first window
tmux rename-window -t $SESSION_NAME:0 "🎯 Project-Lead"

# ============================================================================
# PROJECT MANAGEMENT & ARCHITECTURE
# ============================================================================

# Window 0: Project Lead & Architecture (already created)
tmux send-keys -t $SESSION_NAME:0 "cd $PROJECT_ROOT" C-m
tmux send-keys -t $SESSION_NAME:0 "echo '🎯 PROJECT LEAD & ARCHITECTURE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Role: Overall project coordination, architecture decisions, cross-component integration'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Focus: README.md, docker-compose.yml, Makefile, project-wide decisions'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Key Files: README.md, docker-compose.yml, Makefile, TODO.md, .env.example'" C-m
tmux send-keys -t $SESSION_NAME:0 "echo 'Environment: .env file ready ✅'" C-m
tmux send-keys -t $SESSION_NAME:0 "$AIDER_CMD README.md docker-compose.yml Makefile TODO.md .env.example" C-m

# Window 1: Environment & Configuration Expert
tmux new-window -t $SESSION_NAME -n "🌍 Env-Config" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:1 "echo '🌍 ENVIRONMENT & CONFIGURATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Role: Environment variables, .env management, configuration validation'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Focus: .env, .env.example, environment-specific configurations'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Key Variables: OPENWEBUI_API_KEY, MCP_*, FR0G_*, LOG_LEVEL'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Status: .env file ready ✅'" C-m
tmux send-keys -t $SESSION_NAME:1 "$AIDER_CMD .env.example .env" C-m

# Window 2: Shared Configuration Expert
tmux new-window -t $SESSION_NAME -n "🔧 Config-Expert" -c "$PROJECT_ROOT/pkg/config"
tmux send-keys -t $SESSION_NAME:2 "echo '🔧 SHARED CONFIGURATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Role: Maintain centralized configuration library, validation systems'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Focus: pkg/config/ - validation, loading, shared types'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Key Files: config.go, validation.go, loader.go, examples_test.go'" C-m
tmux send-keys -t $SESSION_NAME:2 "$AIDER_CMD config.go validation.go loader.go examples_test.go README.md" C-m

# ============================================================================
# CORE SUBPROJECT AGENTS
# ============================================================================

# Window 3: fr0g-ai-aip (Core AI Service)
tmux new-window -t $SESSION_NAME -n "🧠 AIP" -c "$PROJECT_ROOT/fr0g-ai-aip"
tmux send-keys -t $SESSION_NAME:3 "echo '🧠 FR0G-AI-AIP AGENT'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Role: Core AI processing engine, persona management, identity processing'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Ports: HTTP :8080, gRPC :9090'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Status: Demographics ✅, Other attributes need implementation'" C-m
tmux send-keys -t $SESSION_NAME:3 "$AIDER_CMD TODO.md internal/" C-m

# Window 4: fr0g-ai-bridge (Integration Service)
tmux new-window -t $SESSION_NAME -n "🌉 Bridge" -c "$PROJECT_ROOT/fr0g-ai-bridge"
tmux send-keys -t $SESSION_NAME:4 "echo '🌉 FR0G-AI-BRIDGE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Role: OpenWebUI integration, REST/gRPC bridge services'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Ports: HTTP :8082, gRPC :9091'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Status: Fully operational ✅'" C-m
tmux send-keys -t $SESSION_NAME:4 "$AIDER_CMD TODO.md internal/" C-m

# Window 5: fr0g-ai-master-control (Cognitive Engine)
tmux new-window -t $SESSION_NAME -n "🎛️ MCP" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:5 "echo '🎛️ FR0G-AI-MASTER-CONTROL AGENT'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Role: System orchestration, cognitive engine, threat processing'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Port: HTTP :8081'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Status: CONSCIOUS AI OPERATIONAL ✅ (Learning Rate: 0.100+)'" C-m
tmux send-keys -t $SESSION_NAME:5 "$AIDER_CMD TODO.md internal/" C-m

# ============================================================================
# ESSENTIAL SUPPORT AGENTS
# ============================================================================

# Window 6: DevOps & Infrastructure
tmux new-window -t $SESSION_NAME -n "🐳 DevOps" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:6 "echo '🐳 DEVOPS & INFRASTRUCTURE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Role: Docker, deployment, CI/CD, infrastructure automation'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Focus: docker-compose.yml, Makefile, .env configuration'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Services: service-registry:8500, aip:8080/9090, bridge:8082/9091, mcp:8081'" C-m
tmux send-keys -t $SESSION_NAME:6 "$AIDER_CMD docker-compose.yml Makefile .env.example" C-m

# Window 7: Build & Test Runner
tmux new-window -t $SESSION_NAME -n "🔨 Build-Test" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:7 "echo '🔨 BUILD & TEST RUNNER'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Role: Build automation, test execution, continuous integration'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Commands: make build-all, make test-all, docker-compose up'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Quick Start: docker-compose up -d'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Available commands:'" C-m
tmux send-keys -t $SESSION_NAME:7 "make help" C-m

# Window 8: System Monitor
tmux new-window -t $SESSION_NAME -n "📈 Monitor" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:8 "echo '📈 SYSTEM MONITOR'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Role: Real-time system monitoring, log viewing, service status'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Commands: docker-compose logs, make health, watch docker ps'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Service Status:'" C-m
tmux send-keys -t $SESSION_NAME:8 "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || echo 'Docker not running - use: docker-compose up -d'" C-m

# Window 9: Interactive Shell
tmux new-window -t $SESSION_NAME -n "💻 Shell" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:9 "echo '💻 INTERACTIVE SHELL'" C-m
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
echo "🎯 AGENT ASSIGNMENTS:"
echo "  0: 🎯 Project-Lead      - Overall coordination & architecture"
echo "  1: 🌍 Env-Config        - Environment & .env management"
echo "  2: 🔧 Config-Expert     - Shared configuration library (pkg/config/)"
echo "  3: 🧠 AIP               - fr0g-ai-aip (Core AI Service) :8080/:9090"
echo "  4: 🌉 Bridge            - fr0g-ai-bridge (Integration) :8082/:9091"
echo "  5: 🎛️ MCP               - fr0g-ai-master-control (Cognitive Engine) :8081"
echo "  6: 🐳 DevOps            - Infrastructure & deployment"
echo "  7: 🔨 Build-Test        - Build & test automation"
echo "  8: 📈 Monitor           - System monitoring"
echo "  9: 💻 Shell             - Interactive shell"
echo ""
echo "🚀 USAGE:"
echo "  tmux attach-session -t $SESSION_NAME"
echo "  Ctrl+b + [0-9]   - Switch to window"
echo "  Ctrl+b + n       - Next window"
echo "  Ctrl+b + p       - Previous window"
echo "  Ctrl+b + w       - Window list"
echo "  Ctrl+b + d       - Detach session"
echo ""
echo "🎯 Each agent window has aider pre-configured for their domain!"
echo "🧠 MCP Intelligence Status: CONSCIOUS AI OPERATIONAL (Learning Rate: 0.100+)"
echo "🔧 Configuration: Centralized system operational (pkg/config/)"
echo "🌍 Environment: .env file ready for configuration"
echo ""
echo "📡 Service Ports:"
echo "  - Service Registry: :8500"
echo "  - AIP Service: HTTP :8080, gRPC :9090"
echo "  - Bridge Service: HTTP :8082, gRPC :9091"
echo "  - Master Control: HTTP :8081"
echo ""
echo "🚀 Quick Start:"
echo "  1. Edit .env file with your API keys and configuration"
echo "  2. Run: docker-compose up -d"
echo "  3. Check health: make health"
echo ""
echo "🔑 Key Environment Variables to Configure:"
echo "  - OPENWEBUI_API_KEY (for OpenWebUI integration)"
echo "  - MCP_LEARNING_ENABLED=true (for AI consciousness)"
echo "  - LOG_LEVEL=info (for appropriate logging)"
echo ""

# Attach to the session
tmux attach-session -t $SESSION_NAME
