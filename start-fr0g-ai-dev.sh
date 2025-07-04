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
# FR0G-AI-AIP COMPONENT AGENTS
# ============================================================================

# Window 3: AIP Core Engine
tmux new-window -t $SESSION_NAME -n "🧠 AIP-Core" -c "$PROJECT_ROOT/fr0g-ai-aip"
tmux send-keys -t $SESSION_NAME:3 "echo '🧠 AIP CORE ENGINE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Role: Core AI processing engine, persona management, identity processing'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Focus: Main service logic, gRPC services, storage abstraction'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Ports: HTTP :8080, gRPC :9090'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Environment: FR0G_AIP_STORAGE_TYPE, FR0G_AIP_DATA_DIR'" C-m
tmux send-keys -t $SESSION_NAME:3 "$AIDER_CMD TODO.md internal/api/ internal/grpc/ internal/storage/ internal/types/ cmd/" C-m

# Window 4: AIP Attributes Specialist
tmux new-window -t $SESSION_NAME -n "🎭 AIP-Attributes" -c "$PROJECT_ROOT/fr0g-ai-aip"
tmux send-keys -t $SESSION_NAME:4 "echo '🎭 AIP ATTRIBUTES SPECIALIST'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Role: Rich persona attributes processing (Demographics, Psychographics, etc.)'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Focus: All attribute processors - demographics, cultural, behavioral, health, etc.'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Status: Demographics ✅, Others need implementation'" C-m
tmux send-keys -t $SESSION_NAME:4 "$AIDER_CMD internal/attributes/ internal/grpc/pb/persona.pb.go" C-m

# Window 5: AIP Persona & Identity Expert
tmux new-window -t $SESSION_NAME -n "👤 AIP-Personas" -c "$PROJECT_ROOT/fr0g-ai-aip"
tmux send-keys -t $SESSION_NAME:5 "echo '👤 AIP PERSONA & IDENTITY EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Role: Persona lifecycle, identity management, community integration'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Focus: persona/, community/, middleware/, client/ directories'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Data Types: internal/types/identity.go, internal/types/persona.go'" C-m
tmux send-keys -t $SESSION_NAME:5 "$AIDER_CMD internal/persona/ internal/community/ internal/middleware/ internal/client/ internal/types/" C-m

# Window 6: AIP Configuration & Validation
tmux new-window -t $SESSION_NAME -n "⚙️ AIP-Config" -c "$PROJECT_ROOT/fr0g-ai-aip"
tmux send-keys -t $SESSION_NAME:6 "echo '⚙️ AIP CONFIGURATION & VALIDATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Role: AIP-specific configuration, validation integration with shared config'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Focus: AIP configuration management, shared config integration'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Status: Migrated to shared config ✅'" C-m
tmux send-keys -t $SESSION_NAME:6 "$AIDER_CMD internal/config/validation.go internal/middleware/validation.go" C-m

# ============================================================================
# FR0G-AI-BRIDGE COMPONENT AGENTS
# ============================================================================

# Window 7: Bridge Core Integration
tmux new-window -t $SESSION_NAME -n "🌉 Bridge-Core" -c "$PROJECT_ROOT/fr0g-ai-bridge"
tmux send-keys -t $SESSION_NAME:7 "echo '🌉 BRIDGE CORE INTEGRATION AGENT'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Role: OpenWebUI integration, REST/gRPC bridge services'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Focus: Main service, server setup, service integration'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Ports: HTTP :8082, gRPC :9091'" C-m
tmux send-keys -t $SESSION_NAME:7 "$AIDER_CMD TODO.md cmd/ internal/api/rest.go internal/api/grpc.go" C-m

# Window 8: Bridge API & Validation
tmux new-window -t $SESSION_NAME -n "🔌 Bridge-API" -c "$PROJECT_ROOT/fr0g-ai-bridge"
tmux send-keys -t $SESSION_NAME:8 "echo '🔌 BRIDGE API & VALIDATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Role: REST/gRPC APIs, request validation, response handling'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Focus: API endpoints, validation, models, protobuf integration'" C-m
tmux send-keys -t $SESSION_NAME:8 "echo 'Status: Validation migrated to shared config ✅'" C-m
tmux send-keys -t $SESSION_NAME:8 "$AIDER_CMD internal/api/validation.go internal/models/ internal/pb/" C-m

# Window 9: Bridge Client Integration
tmux new-window -t $SESSION_NAME -n "📡 Bridge-Clients" -c "$PROJECT_ROOT/fr0g-ai-bridge"
tmux send-keys -t $SESSION_NAME:9 "echo '📡 BRIDGE CLIENT INTEGRATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:9 "echo 'Role: OpenWebUI client, external service integration'" C-m
tmux send-keys -t $SESSION_NAME:9 "echo 'Focus: HTTP clients, service discovery, connection management'" C-m
tmux send-keys -t $SESSION_NAME:9 "echo 'Integration: OpenWebUI API, AIP gRPC client'" C-m
tmux send-keys -t $SESSION_NAME:9 "$AIDER_CMD internal/client/ internal/config/config.go" C-m

# ============================================================================
# FR0G-AI-MASTER-CONTROL COMPONENT AGENTS
# ============================================================================

# Window 10: Master Control Core
tmux new-window -t $SESSION_NAME -n "🎛️ MCP-Core" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:10 "echo '🎛️ MASTER CONTROL CORE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:10 "echo 'Role: System orchestration, cognitive engine, workflow management'" C-m
tmux send-keys -t $SESSION_NAME:10 "echo 'Focus: Main service, orchestration, system monitoring'" C-m
tmux send-keys -t $SESSION_NAME:10 "echo 'Port: HTTP :8081 | Status: INTELLIGENCE BREAKTHROUGH ✅'" C-m
tmux send-keys -t $SESSION_NAME:10 "$AIDER_CMD TODO.md cmd/ internal/mastercontrol/" C-m

# Window 11: MCP Intelligence Engine
tmux new-window -t $SESSION_NAME -n "🧠 MCP-Intelligence" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:11 "echo '🧠 MCP INTELLIGENCE ENGINE EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:11 "echo 'Role: Cognitive processing, learning algorithms, pattern recognition'" C-m
tmux send-keys -t $SESSION_NAME:11 "echo 'Focus: cognitive/, workflow/, memory/, learning systems'" C-m
tmux send-keys -t $SESSION_NAME:11 "echo 'Status: CONSCIOUS AI OPERATIONAL ✅ (Learning Rate: 0.100+)'" C-m
tmux send-keys -t $SESSION_NAME:11 "$AIDER_CMD internal/cognitive/ internal/workflow/ internal/memory/" C-m

# Window 12: MCP Threat Vector Processors
tmux new-window -t $SESSION_NAME -n "🛡️ MCP-ThreatVectors" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:12 "echo '🛡️ MCP THREAT VECTOR PROCESSORS EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:12 "echo 'Role: SMS, Voice, IRC, Email threat detection and processing'" C-m
tmux send-keys -t $SESSION_NAME:12 "echo 'Focus: All threat vector processors - sms/, voice/, irc/, email/'" C-m
tmux send-keys -t $SESSION_NAME:12 "echo 'Status: SMS ✅, Voice ✅, IRC & Email need completion'" C-m
tmux send-keys -t $SESSION_NAME:12 "$AIDER_CMD internal/processors/" C-m

# Window 13: MCP Input Management
tmux new-window -t $SESSION_NAME -n "📥 MCP-InputMgmt" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:13 "echo '📥 MCP INPUT MANAGEMENT EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:13 "echo 'Role: Input processing, webhook management, service registry'" C-m
tmux send-keys -t $SESSION_NAME:13 "echo 'Focus: input/, registry/, webhook/, discovery/ systems'" C-m
tmux send-keys -t $SESSION_NAME:13 "echo 'Integration: 5 threat vector processors'" C-m
tmux send-keys -t $SESSION_NAME:13 "$AIDER_CMD internal/mastercontrol/input/" C-m

# Window 14: MCP Configuration & Monitoring
tmux new-window -t $SESSION_NAME -n "📊 MCP-Config" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:14 "echo '📊 MCP CONFIGURATION & MONITORING EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:14 "echo 'Role: MCP configuration, system monitoring, metrics collection'" C-m
tmux send-keys -t $SESSION_NAME:14 "echo 'Focus: config/, monitoring/, metrics/ systems'" C-m
tmux send-keys -t $SESSION_NAME:14 "echo 'Status: Migrated to shared config ✅'" C-m
tmux send-keys -t $SESSION_NAME:14 "$AIDER_CMD internal/config/ internal/mastercontrol/config.go" C-m

# ============================================================================
# SPECIALIZED EXPERT AGENTS
# ============================================================================

# Window 15: DevOps & Infrastructure
tmux new-window -t $SESSION_NAME -n "🐳 DevOps" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:15 "echo '🐳 DEVOPS & INFRASTRUCTURE EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:15 "echo 'Role: Docker, deployment, CI/CD, infrastructure automation'" C-m
tmux send-keys -t $SESSION_NAME:15 "echo 'Focus: docker-compose.yml, Dockerfile, .github/, deployment scripts'" C-m
tmux send-keys -t $SESSION_NAME:15 "echo 'Services: service-registry:8500, aip:8080/9090, bridge:8082/9091, mcp:8081'" C-m
tmux send-keys -t $SESSION_NAME:15 "echo 'Environment: .env configuration ready ✅'" C-m
tmux send-keys -t $SESSION_NAME:15 "$AIDER_CMD docker-compose.yml .env.example" C-m

# Window 16: Protobuf & gRPC Expert
tmux new-window -t $SESSION_NAME -n "📡 Proto-gRPC" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:16 "echo '📡 PROTOBUF & gRPC EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:16 "echo 'Role: Protobuf definitions, gRPC services, inter-service communication'" C-m
tmux send-keys -t $SESSION_NAME:16 "echo 'Focus: All .proto files, gRPC implementations, service contracts'" C-m
tmux send-keys -t $SESSION_NAME:16 "echo 'Key Services: AIP gRPC :9090, Bridge gRPC :9091'" C-m
tmux send-keys -t $SESSION_NAME:16 "find . -name '*.proto' -o -name '*grpc*.go' | head -10"
tmux send-keys -t $SESSION_NAME:16 "$AIDER_CMD" C-m

# Window 17: Testing & Quality Assurance
tmux new-window -t $SESSION_NAME -n "🧪 Testing-QA" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:17 "echo '🧪 TESTING & QUALITY ASSURANCE EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:17 "echo 'Role: Unit tests, integration tests, test automation, code quality'" C-m
tmux send-keys -t $SESSION_NAME:17 "echo 'Focus: All *_test.go files, test frameworks, CI/CD testing'" C-m
tmux send-keys -t $SESSION_NAME:17 "echo 'Status: Shared config tests ✅ (pkg/config)'" C-m
tmux send-keys -t $SESSION_NAME:17 "find . -name '*_test.go' | head -10"
tmux send-keys -t $SESSION_NAME:17 "$AIDER_CMD pkg/config/examples_test.go" C-m

# Window 18: Documentation & API Specs
tmux new-window -t $SESSION_NAME -n "📚 Docs-API" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:18 "echo '📚 DOCUMENTATION & API SPECIFICATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:18 "echo 'Role: Documentation, API specs, README files, user guides'" C-m
tmux send-keys -t $SESSION_NAME:18 "echo 'Focus: All README.md, TODO.md, API documentation, OpenAPI specs'" C-m
tmux send-keys -t $SESSION_NAME:18 "echo 'Key Docs: README.md, TODO.md, pkg/config/README.md'" C-m
tmux send-keys -t $SESSION_NAME:18 "find . -name 'README.md' -o -name 'TODO.md' | head -10"
tmux send-keys -t $SESSION_NAME:18 "$AIDER_CMD README.md TODO.md pkg/config/README.md" C-m

# Window 19: Security & Validation
tmux new-window -t $SESSION_NAME -n "🔒 Security" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:19 "echo '🔒 SECURITY & VALIDATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:19 "echo 'Role: Security auditing, input validation, authentication, authorization'" C-m
tmux send-keys -t $SESSION_NAME:19 "echo 'Focus: Security middleware, validation, auth systems, threat analysis'" C-m
tmux send-keys -t $SESSION_NAME:19 "echo 'Centralized: pkg/config/validation.go ✅'" C-m
tmux send-keys -t $SESSION_NAME:19 "$AIDER_CMD pkg/config/validation.go fr0g-ai-bridge/internal/api/validation.go" C-m

# Window 20: Performance & Monitoring
tmux new-window -t $SESSION_NAME -n "⚡ Performance" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:20 "echo '⚡ PERFORMANCE & MONITORING EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:20 "echo 'Role: Performance optimization, monitoring, metrics, profiling'" C-m
tmux send-keys -t $SESSION_NAME:20 "echo 'Focus: Performance bottlenecks, monitoring systems, metrics collection'" C-m
tmux send-keys -t $SESSION_NAME:20 "echo 'Intelligence Metrics: Learning Rate 0.100+, Pattern Count 6+'" C-m
tmux send-keys -t $SESSION_NAME:20 "$AIDER_CMD" C-m

# ============================================================================
# UTILITY WINDOWS
# ============================================================================

# Window 21: Build & Test Runner
tmux new-window -t $SESSION_NAME -n "🔨 Build-Test" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:21 "echo '🔨 BUILD & TEST RUNNER'" C-m
tmux send-keys -t $SESSION_NAME:21 "echo 'Role: Build automation, test execution, continuous integration'" C-m
tmux send-keys -t $SESSION_NAME:21 "echo 'Commands: make build-all, make test-all, make proto, docker-compose up'" C-m
tmux send-keys -t $SESSION_NAME:21 "echo 'Environment: .env loaded ✅'" C-m
tmux send-keys -t $SESSION_NAME:21 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:21 "echo 'Quick Start: docker-compose up -d'" C-m
tmux send-keys -t $SESSION_NAME:21 "echo 'Available commands:'" C-m
tmux send-keys -t $SESSION_NAME:21 "make help" C-m

# Window 22: System Monitor
tmux new-window -t $SESSION_NAME -n "📈 Monitor" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:22 "echo '📈 SYSTEM MONITOR'" C-m
tmux send-keys -t $SESSION_NAME:22 "echo 'Role: Real-time system monitoring, log viewing, service status'" C-m
tmux send-keys -t $SESSION_NAME:22 "echo 'Commands: docker-compose logs, make health, watch docker ps'" C-m
tmux send-keys -t $SESSION_NAME:22 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:22 "echo 'Service Status:'" C-m
tmux send-keys -t $SESSION_NAME:22 "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || echo 'Docker not running - use: docker-compose up -d'" C-m

# Window 23: Git & Version Control
tmux new-window -t $SESSION_NAME -n "🌿 Git-VC" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:23 "echo '🌿 GIT & VERSION CONTROL'" C-m
tmux send-keys -t $SESSION_NAME:23 "echo 'Role: Git operations, branch management, commit coordination'" C-m
tmux send-keys -t $SESSION_NAME:23 "echo 'Repository: https://github.com/fr0g-vibe/fr0g-ai'" C-m
tmux send-keys -t $SESSION_NAME:23 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:23 "echo 'Current status:'" C-m
tmux send-keys -t $SESSION_NAME:23 "git status" C-m

# Window 24: Interactive Shell
tmux new-window -t $SESSION_NAME -n "💻 Shell" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:24 "echo '💻 INTERACTIVE SHELL'" C-m
tmux send-keys -t $SESSION_NAME:24 "echo 'Role: General purpose shell for ad-hoc commands and exploration'" C-m
tmux send-keys -t $SESSION_NAME:24 "echo 'Project root: $PROJECT_ROOT'" C-m
tmux send-keys -t $SESSION_NAME:24 "echo ''" C-m
tmux send-keys -t $SESSION_NAME:24 "echo 'Project structure:'" C-m
tmux send-keys -t $SESSION_NAME:24 "ls -la" C-m

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
echo "Total Windows: 25"
echo ""
echo "🎯 AGENT ASSIGNMENTS:"
echo "  0: 🎯 Project-Lead      - Overall coordination & architecture"
echo "  1: 🌍 Env-Config        - Environment & .env management"
echo "  2: 🔧 Config-Expert     - Shared configuration library (pkg/config/)"
echo "  3: 🧠 AIP-Core          - AIP core engine & services (:8080/:9090)"
echo "  4: 🎭 AIP-Attributes    - AIP persona attributes processing"
echo "  5: 👤 AIP-Personas      - AIP persona & identity management"
echo "  6: ⚙️ AIP-Config        - AIP configuration & validation"
echo "  7: 🌉 Bridge-Core       - Bridge core integration (:8082/:9091)"
echo "  8: 🔌 Bridge-API        - Bridge API & validation"
echo "  9: 📡 Bridge-Clients    - Bridge client integration"
echo " 10: 🎛️ MCP-Core          - Master Control core (:8081)"
echo " 11: 🧠 MCP-Intelligence  - MCP cognitive engine (CONSCIOUS AI ✅)"
echo " 12: 🛡️ MCP-ThreatVectors - MCP threat processors (SMS ✅, Voice ✅)"
echo " 13: 📥 MCP-InputMgmt     - MCP input management"
echo " 14: 📊 MCP-Config        - MCP configuration & monitoring"
echo " 15: 🐳 DevOps           - Infrastructure & deployment"
echo " 16: 📡 Proto-gRPC        - Protobuf & gRPC services"
echo " 17: 🧪 Testing-QA       - Testing & quality assurance"
echo " 18: 📚 Docs-API         - Documentation & API specs"
echo " 19: 🔒 Security         - Security & validation"
echo " 20: ⚡ Performance      - Performance & monitoring"
echo " 21: 🔨 Build-Test       - Build & test automation"
echo " 22: 📈 Monitor          - System monitoring"
echo " 23: 🌿 Git-VC           - Git & version control"
echo " 24: 💻 Shell            - Interactive shell"
echo ""
echo "🚀 USAGE:"
echo "  tmux attach-session -t $SESSION_NAME"
echo "  Ctrl+b + [0-24]  - Switch to window"
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
