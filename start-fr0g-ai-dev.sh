#!/bin/bash

# fr0g-ai Development Environment Setup Script
# Creates tmux session with dedicated windows for each coding agent/component

set -e

# Configuration
SESSION_NAME="fr0g-ai"
PROJECT_ROOT="$(pwd)"
AIDER_CMD="aider --dark-mode"

# Function to create aider command - just start aider with files
create_aider_cmd() {
    local prompt_file="$1"
    local files="$2"
    echo "$AIDER_CMD $files"
}

# Function to send system prompt after aider starts
send_system_prompt() {
    local window="$1"
    local prompt_file="$2"
    
    # Check if prompt file exists before trying to use it
    if [[ ! -f "$prompt_file" ]]; then
        echo "WARNING: Prompt file $prompt_file not found, skipping system prompt for window $window"
        return 0
    fi
    
    # Wait for aider to start and be ready
    sleep 8
    
    # Check if the tmux window still exists
    if ! tmux list-windows -t $SESSION_NAME 2>/dev/null | grep -q "^$window:"; then
        echo "WARNING: Window $window no longer exists, skipping system prompt"
        return 0
    fi
    
    # Read the prompt file content
    local prompt_content
    if ! prompt_content=$(cat "$prompt_file" 2>/dev/null); then
        echo "WARNING: Could not read prompt file $prompt_file"
        return 0
    fi
    
    # Send the system prompt as a regular chat message
    if [[ -n "$prompt_content" ]]; then
        # Use a simpler approach - send the entire prompt at once using tmux's literal mode
        # First, create a temporary file with the prompt content
        local temp_file=$(mktemp)
        echo "$prompt_content" > "$temp_file"
        
        # Send the content using tmux load-buffer and paste-buffer
        tmux load-buffer "$temp_file"
        tmux paste-buffer -t $SESSION_NAME:$window
        tmux send-keys -t $SESSION_NAME:$window C-m
        
        # Clean up temp file
        rm "$temp_file"
    fi
}


echo "STARTING Starting fr0g-ai Development Environment..."
echo "📁 Project Root: $PROJECT_ROOT"

# Verify we're in the right directory
if [[ ! -f "docker-compose.yml" ]] || [[ ! -f "Makefile" ]] || [[ ! -f ".env.example" ]]; then
    echo "FAILED Error: Please run this script from the fr0g-ai project root directory"
    echo "   Expected files: docker-compose.yml, Makefile, .env.example"
    exit 1
fi

# Check if .env exists, if not copy from .env.example
if [[ ! -f ".env" ]]; then
    echo "WARNING  No .env file found. Creating from .env.example..."
    cp .env.example .env
    echo "COMPLETED Created .env file from .env.example"
    echo "NOTES Please edit .env file with your specific configuration"
else
    echo "COMPLETED Existing .env file found - preserving your configuration"
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
tmux send-keys -t $SESSION_NAME:0 "echo 'Environment: .env file ready COMPLETED'" C-m
LEAD_CMD=$(create_aider_cmd ".prompts/project-lead.md" "--file README.md --file docker-compose.yml --file Makefile --file TODO.md --file .env.example")
tmux send-keys -t $SESSION_NAME:0 "$LEAD_CMD" C-m

# ============================================================================
# CORE SUBPROJECT AGENTS
# ============================================================================

# Window 1: fr0g-ai-aip (Core AI Service)
tmux new-window -t $SESSION_NAME -n "AIP" -c "$PROJECT_ROOT/fr0g-ai-aip"
tmux send-keys -t $SESSION_NAME:1 "echo 'FR0G-AI-AIP AGENT'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Role: Core AI processing engine, persona management, identity processing'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Ports: HTTP :8080, gRPC :9090'" C-m
tmux send-keys -t $SESSION_NAME:1 "echo 'Status: Demographics COMPLETED, Other attributes need implementation'" C-m
tmux send-keys -t $SESSION_NAME:1 "cd fr0g-ai-aip" C-m
AIP_CMD=$(create_aider_cmd ".prompts/aip-agent.md" "--file TODO.md")
tmux send-keys -t $SESSION_NAME:1 "$AIP_CMD" C-m

# Window 2: fr0g-ai-bridge (Integration Service)
tmux new-window -t $SESSION_NAME -n "Bridge" -c "$PROJECT_ROOT/fr0g-ai-bridge"
tmux send-keys -t $SESSION_NAME:2 "echo 'FR0G-AI-BRIDGE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Role: OpenWebUI integration, REST/gRPC bridge services'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Ports: HTTP :8082, gRPC :9091'" C-m
tmux send-keys -t $SESSION_NAME:2 "echo 'Status: Fully operational COMPLETED'" C-m
tmux send-keys -t $SESSION_NAME:2 "cd fr0g-ai-bridge" C-m
BRIDGE_CMD=$(create_aider_cmd ".prompts/bridge-agent.md" "--file TODO.md")
tmux send-keys -t $SESSION_NAME:2 "$BRIDGE_CMD" C-m

# Window 3: fr0g-ai-master-control (Cognitive Engine)
tmux new-window -t $SESSION_NAME -n "MCP" -c "$PROJECT_ROOT/fr0g-ai-master-control"
tmux send-keys -t $SESSION_NAME:3 "echo 'FR0G-AI-MASTER-CONTROL AGENT'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Role: System orchestration, cognitive engine, threat processing'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Port: HTTP :8081'" C-m
tmux send-keys -t $SESSION_NAME:3 "echo 'Status: CONSCIOUS AI OPERATIONAL COMPLETED (Learning Rate: 0.100+)'" C-m
tmux send-keys -t $SESSION_NAME:3 "cd fr0g-ai-master-control" C-m
MCP_CMD=$(create_aider_cmd ".prompts/mcp-agent.md" "--file TODO.md")
tmux send-keys -t $SESSION_NAME:3 "$MCP_CMD" C-m

# Window 4: fr0g-ai-io (Input/Output Processing)
tmux new-window -t $SESSION_NAME -n "IO" -c "$PROJECT_ROOT/fr0g-ai-io"
tmux send-keys -t $SESSION_NAME:4 "echo 'FR0G-AI-IO AGENT'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Role: Input/Output processing, threat vector handling, external integrations'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Ports: HTTP :8083, gRPC :9092'" C-m
tmux send-keys -t $SESSION_NAME:4 "echo 'Status: NEW SERVICE - SMS, Voice, IRC, ESMTP, Discord processors'" C-m
tmux send-keys -t $SESSION_NAME:4 "cd fr0g-ai-io" C-m
IO_CMD=$(create_aider_cmd ".prompts/io-agent.md" "--file TODO.md")
tmux send-keys -t $SESSION_NAME:4 "$IO_CMD" C-m

# Window 5: Configuration Expert
tmux new-window -t $SESSION_NAME -n "Config" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:5 "echo 'CONFIGURATION EXPERT'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Role: Environment variables, shared config library, validation systems'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Focus: .env, pkg/config/ - validation, loading, shared types'" C-m
tmux send-keys -t $SESSION_NAME:5 "echo 'Key Files: .env, pkg/config/*.go'" C-m
CONFIG_CMD=$(create_aider_cmd ".prompts/config-agent.md" "--file .env.example --file .env --file pkg/config/config.go --file pkg/config/validation.go --file pkg/config/loader.go --file pkg/config/examples_test.go --file pkg/config/README.md")
tmux send-keys -t $SESSION_NAME:5 "$CONFIG_CMD" C-m

# ============================================================================
# ESSENTIAL SUPPORT AGENTS
# ============================================================================

# Window 6: DevOps & Infrastructure
tmux new-window -t $SESSION_NAME -n "DevOps" -c "$PROJECT_ROOT"
tmux send-keys -t $SESSION_NAME:6 "echo 'DEVOPS & INFRASTRUCTURE AGENT'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Role: Docker, deployment, CI/CD, infrastructure automation'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Focus: docker-compose.yml, Makefile, .env configuration'" C-m
tmux send-keys -t $SESSION_NAME:6 "echo 'Services: service-registry:8500, aip:8080/9090, bridge:8082/9091, mcp:8081, io:8083/9092'" C-m
DEVOPS_CMD=$(create_aider_cmd ".prompts/devops-agent.md" "--file docker-compose.yml --file Makefile --file .env.example")
tmux send-keys -t $SESSION_NAME:6 "$DEVOPS_CMD" C-m

# Window 7: Registry Agent
tmux new-window -t $SESSION_NAME -n "Registry" -c "$PROJECT_ROOT/fr0g-ai-registry"
tmux send-keys -t $SESSION_NAME:7 "echo 'FR0G-AI-REGISTRY AGENT'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Role: Service discovery, registration, and health monitoring'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Port: HTTP :8500'" C-m
tmux send-keys -t $SESSION_NAME:7 "echo 'Status: Extracted from master-control - ready for enhancement'" C-m
tmux send-keys -t $SESSION_NAME:7 "cd fr0g-ai-registry" C-m
REGISTRY_CMD=$(create_aider_cmd ".prompts/registry-agent.md" "--file TODO.md")
tmux send-keys -t $SESSION_NAME:7 "$REGISTRY_CMD" C-m

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
# SYSTEM PROMPT SETUP
# ============================================================================

echo "PROCESSING Setting up system prompts for specialized agents..."

# Wait for all aider instances to start
sleep 3

# Send system prompts (with correct paths from project root)
{
    send_system_prompt 0 ".prompts/project-lead.md"
    send_system_prompt 1 ".prompts/aip-agent.md" 
    send_system_prompt 2 ".prompts/bridge-agent.md"
    send_system_prompt 3 ".prompts/mcp-agent.md"
    send_system_prompt 4 ".prompts/io-agent.md"
    send_system_prompt 5 ".prompts/config-agent.md"
    send_system_prompt 6 ".prompts/devops-agent.md"
    send_system_prompt 7 ".prompts/registry-agent.md"
} &

# Store the background job PID
PROMPT_SETUP_PID=$!

echo "COMPLETED System prompts being configured in background..."

# ============================================================================
# SESSION FINALIZATION
# ============================================================================

# Add cleanup function
cleanup() {
    echo "Cleaning up background processes..."
    if [[ -n "$PROMPT_SETUP_PID" ]]; then
        kill $PROMPT_SETUP_PID 2>/dev/null || true
    fi
    jobs -p | xargs -r kill 2>/dev/null || true
}

# Set trap for cleanup
trap cleanup EXIT

# Return to first window (Project Lead)
tmux select-window -t $SESSION_NAME:0

# Display session info
echo ""
echo "COMPLETED fr0g-ai Development Environment Created!"
echo ""
echo "LIST TMUX SESSION OVERVIEW:"
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
echo "STARTING USAGE:"
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
echo "🤖 AGENT DISPATCH SYSTEM:"
echo "  Project Lead can dispatch commands to specialized agents:"
echo "  tmux send-keys -t fr0g-ai:1 \"Implement persona CRUD operations\" C-m"
echo "  tmux send-keys -t fr0g-ai:2 \"Add health check validation\" C-m"
echo "  tmux send-keys -t fr0g-ai:3 \"Optimize learning algorithms\" C-m"
echo "  tmux send-keys -t fr0g-ai:8 \"make build-all\" C-m"
echo ""
echo "📋 AGENT SPECIALIZATIONS:"
echo "  0: Project-Lead    - Architecture & coordination"
echo "  1: AIP Agent       - Core AI processing engine"
echo "  2: Bridge Agent    - External integrations & API gateway"
echo "  3: MCP Agent       - Cognitive intelligence engine"
echo "  4: IO Agent        - Input/output & threat processing"
echo "  5: Config Agent    - Configuration & environment mgmt"
echo "  6: DevOps Agent    - Infrastructure & deployment"
echo "  7: Registry Agent  - Service discovery & health monitor"
echo "  8: Build-Test      - Build automation & testing"
echo "  9: Shell           - Interactive shell & ad-hoc commands"
echo ""
echo "Quick Start:"
echo "  1. Edit .env file with your API keys and configuration"
echo "  2. Run: docker-compose up -d"
echo "  3. Check health: make health"
echo "  4. Use dispatch system for coordinated development"
echo ""
echo "Key Environment Variables to Configure:"
echo "  - OPENWEBUI_API_KEY (for OpenWebUI integration)"
echo "  - MCP_LEARNING_ENABLED=true (for AI consciousness)"
echo "  - LOG_LEVEL=info (for appropriate logging)"
echo ""
echo "COMPLETED fr0g-ai Development Environment Ready!"

# Wait a moment for background setup to complete
echo "Waiting for agent setup to complete..."
sleep 5

# Attach to the session with error handling
echo "Attaching to tmux session $SESSION_NAME..."
if ! tmux attach-session -t $SESSION_NAME; then
    echo "ERROR: Failed to attach to tmux session"
    echo "Session created but attachment failed"
    echo "Try manually: tmux attach-session -t $SESSION_NAME"
    exit 1
fi
