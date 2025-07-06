# Agent Dispatch System Documentation

## Overview

The fr0g.ai Agent Dispatch System is a sophisticated multi-agent development environment that enables coordinated development across all components using tmux and aider. This system allows the Project Lead to dispatch specific tasks to specialized AI agents, each with domain expertise and dedicated system prompts.

## Architecture

### Agent Window Layout

```
┌─────────────────────────────────────────────────────────────────┐
│                    fr0g-ai Development Environment               │
├─────────────────────────────────────────────────────────────────┤
│ Window 0: Project-Lead    │ Architecture & Cross-Component      │
│ Window 1: AIP Agent       │ Core AI Processing Engine           │
│ Window 2: Bridge Agent    │ External Integrations & API Gateway │
│ Window 3: MCP Agent       │ Cognitive Intelligence Engine       │
│ Window 4: IO Agent        │ Input/Output & Threat Processing    │
│ Window 5: Config Agent    │ Configuration & Environment Mgmt    │
│ Window 6: DevOps Agent    │ Infrastructure & Deployment         │
│ Window 7: Registry Agent  │ Service Discovery & Health Monitor  │
│ Window 8: Build-Test      │ Build Automation & Testing          │
│ Window 9: Shell           │ Interactive Shell & Ad-hoc Commands │
└─────────────────────────────────────────────────────────────────┘
```

### Agent Specializations

#### Project Lead (Window 0)
- **Domain**: Overall project coordination, architecture decisions, cross-component integration
- **Files**: README.md, docker-compose.yml, Makefile, TODO.md, .env.example
- **Responsibilities**: 
  - Coordinate cross-component changes
  - Maintain service port assignments
  - Ensure consistent environment variable naming
  - Review architectural decisions
  - Dispatch tasks to specialized agents

#### AIP Agent (Window 1)
- **Domain**: Core AI processing, persona management, identity processing
- **Component**: fr0g-ai-aip
- **Ports**: HTTP :8080, gRPC :9090
- **Responsibilities**:
  - Persona CRUD operations and management
  - Rich attribute processing (Demographics, Psychographics, etc.)
  - gRPC service implementation and optimization
  - Database migration planning
  - Performance optimization for 1000+ concurrent users

#### Bridge Agent (Window 2)
- **Domain**: External integrations, API gateway, protocol translation
- **Component**: fr0g-ai-bridge
- **Ports**: HTTP :8082, gRPC :9091
- **Responsibilities**:
  - OpenWebUI API integration and optimization
  - Multi-LLM provider support (OpenAI, Anthropic, Cohere)
  - API gateway functionality and routing
  - Authentication and authorization systems
  - Rate limiting and quota management

#### MCP Agent (Window 3)
- **Domain**: Cognitive intelligence, orchestration, conscious AI
- **Component**: fr0g-ai-master-control
- **Ports**: HTTP :8081
- **Responsibilities**:
  - Cognitive intelligence engine optimization
  - Adaptive learning algorithms and neural adaptation
  - Workflow orchestration and automation
  - Threat analysis and pattern recognition
  - Memory management and knowledge graphs

#### IO Agent (Window 4)
- **Domain**: Input/output processing, threat vector handling, external integrations
- **Component**: fr0g-ai-io
- **Ports**: HTTP :8083, gRPC :9092
- **Responsibilities**:
  - Threat vector interception and analysis
  - Input processor optimization (SMS, Voice, IRC, ESMTP, Discord)
  - Output command validation and delivery
  - External API integration (Google Voice, etc.)
  - Real-time communication monitoring

#### Config Agent (Window 5)
- **Domain**: Configuration and environment specialist
- **Files**: .env, .env.example, pkg/config/*.go
- **Responsibilities**:
  - Maintain centralized configuration system in pkg/config/
  - Ensure all components use shared config library
  - Validate environment variable consistency
  - Implement new validation functions as needed
  - Maintain .env.example with all required variables

#### DevOps Agent (Window 6)
- **Domain**: Infrastructure and deployment specialist
- **Files**: docker-compose.yml, Makefile, Dockerfiles, deployment scripts
- **Responsibilities**:
  - Maintain Docker Compose orchestration
  - Optimize Dockerfile builds and security
  - Manage service port assignments and networking
  - Implement CI/CD pipelines and automation
  - Monitor infrastructure health and performance

#### Registry Agent (Window 7)
- **Domain**: Service discovery, registration, health monitoring
- **Component**: fr0g-ai-registry
- **Ports**: HTTP :8500
- **Responsibilities**:
  - Service discovery and registration optimization
  - Health monitoring and automated checks
  - Load balancing and service routing
  - Performance optimization (target <5ms discovery)
  - Redis persistence for zero data loss

## Dispatch Commands

### Basic Syntax

```bash
# Command Format:
tmux send-keys -t fr0g-ai:WINDOW_NUMBER "COMMAND" C-m

# Windows 0-7: Aider AI Agents (with specialized system prompts)
# Windows 8-9: Shell environments (direct command execution)
```

### Agent-Specific Dispatch Examples

#### Core AI Development (AIP Agent - Window 1)
```bash
# Persona management
tmux send-keys -t fr0g-ai:1 "Implement persona CRUD operations with comprehensive validation" C-m
tmux send-keys -t fr0g-ai:1 "Add rich attribute processors for Demographics and Psychographics" C-m
tmux send-keys -t fr0g-ai:1 "Create persona recommendation engine with similarity algorithms" C-m

# Performance optimization
tmux send-keys -t fr0g-ai:1 "Optimize gRPC service performance for 1000+ concurrent users" C-m
tmux send-keys -t fr0g-ai:1 "Implement caching layer for high-frequency persona lookups" C-m

# Database migration
tmux send-keys -t fr0g-ai:1 "Migrate file storage to PostgreSQL with connection pooling" C-m
tmux send-keys -t fr0g-ai:1 "Add database transaction support for persona operations" C-m
```

#### Integration Development (Bridge Agent - Window 2)
```bash
# Multi-LLM support
tmux send-keys -t fr0g-ai:2 "Add multi-LLM provider support for OpenAI and Anthropic" C-m
tmux send-keys -t fr0g-ai:2 "Implement intelligent model routing based on request type" C-m

# API management
tmux send-keys -t fr0g-ai:2 "Implement API rate limiting and quota management" C-m
tmux send-keys -t fr0g-ai:2 "Add comprehensive health check validation with metrics" C-m
tmux send-keys -t fr0g-ai:2 "Enhance OpenWebUI integration with streaming responses" C-m

# Security
tmux send-keys -t fr0g-ai:2 "Implement OAuth2 and JWT authentication" C-m
tmux send-keys -t fr0g-ai:2 "Add API key management and rotation" C-m
```

#### Cognitive Engine Development (MCP Agent - Window 3)
```bash
# Learning algorithms
tmux send-keys -t fr0g-ai:3 "Optimize learning rate algorithms for threat detection" C-m
tmux send-keys -t fr0g-ai:3 "Implement advanced neural adaptation mechanisms" C-m

# Autonomous capabilities
tmux send-keys -t fr0g-ai:3 "Implement autonomous workflow generation capabilities" C-m
tmux send-keys -t fr0g-ai:3 "Add self-healing and auto-remediation features" C-m

# Threat analysis
tmux send-keys -t fr0g-ai:3 "Add predictive threat modeling with 24-hour forecasting" C-m
tmux send-keys -t fr0g-ai:3 "Enhance cognitive reflection cycles for better adaptation" C-m
```

#### I/O Processing Development (IO Agent - Window 4)
```bash
# External integrations
tmux send-keys -t fr0g-ai:4 "Complete SMS processor integration with Google Voice API" C-m
tmux send-keys -t fr0g-ai:4 "Implement voice processing with speech-to-text integration" C-m

# Threat detection
tmux send-keys -t fr0g-ai:4 "Implement real-time threat detection for IRC channels" C-m
tmux send-keys -t fr0g-ai:4 "Add automated response generation for email threats" C-m
tmux send-keys -t fr0g-ai:4 "Optimize ESMTP processor for high-volume email analysis" C-m

# Communication monitoring
tmux send-keys -t fr0g-ai:4 "Add multi-channel output coordination" C-m
tmux send-keys -t fr0g-ai:4 "Implement communication pattern analysis" C-m
```

#### Configuration Management (Config Agent - Window 5)
```bash
# Hot-reload capabilities
tmux send-keys -t fr0g-ai:5 "Add hot-reload capabilities for configuration changes" C-m
tmux send-keys -t fr0g-ai:5 "Implement configuration change notifications" C-m

# Validation
tmux send-keys -t fr0g-ai:5 "Implement advanced validation rules for security configs" C-m
tmux send-keys -t fr0g-ai:5 "Add cross-component configuration consistency checks" C-m

# Templates and generators
tmux send-keys -t fr0g-ai:5 "Create configuration templates for different environments" C-m
tmux send-keys -t fr0g-ai:5 "Add configuration audit and compliance checking" C-m
```

#### Infrastructure Development (DevOps Agent - Window 6)
```bash
# Container security
tmux send-keys -t fr0g-ai:6 "Implement production-ready container security hardening" C-m
tmux send-keys -t fr0g-ai:6 "Add container vulnerability scanning" C-m

# Monitoring
tmux send-keys -t fr0g-ai:6 "Add Prometheus metrics and Grafana dashboards" C-m
tmux send-keys -t fr0g-ai:6 "Implement distributed tracing with correlation IDs" C-m

# CI/CD
tmux send-keys -t fr0g-ai:6 "Create CI/CD pipeline with automated testing" C-m
tmux send-keys -t fr0g-ai:6 "Optimize Docker builds for faster deployment cycles" C-m
```

#### Service Discovery Development (Registry Agent - Window 7)
```bash
# Performance optimization
tmux send-keys -t fr0g-ai:7 "Optimize service discovery performance to <5ms latency" C-m
tmux send-keys -t fr0g-ai:7 "Implement Redis persistence for zero data loss" C-m

# Health monitoring
tmux send-keys -t fr0g-ai:7 "Add automated health checking with configurable intervals" C-m
tmux send-keys -t fr0g-ai:7 "Implement service mesh integration" C-m

# Load balancing
tmux send-keys -t fr0g-ai:7 "Enhance service registry API with load balancing support" C-m
tmux send-keys -t fr0g-ai:7 "Add weighted round-robin routing" C-m
```

### Shell Command Dispatch

#### Build and Test Automation (Window 8)
```bash
# Build operations
tmux send-keys -t fr0g-ai:8 "make build-all" C-m
tmux send-keys -t fr0g-ai:8 "make clean-all && make build-all" C-m

# Testing
tmux send-keys -t fr0g-ai:8 "make test-all-integration" C-m
tmux send-keys -t fr0g-ai:8 "make test-performance" C-m
tmux send-keys -t fr0g-ai:8 "make validate-production" C-m

# Docker operations
tmux send-keys -t fr0g-ai:8 "docker-compose up -d" C-m
tmux send-keys -t fr0g-ai:8 "docker-compose down && docker-compose up -d" C-m

# Health checks
tmux send-keys -t fr0g-ai:8 "make health" C-m
tmux send-keys -t fr0g-ai:8 "make health" C-m
```

#### General Shell Operations (Window 9)
```bash
# Git operations
tmux send-keys -t fr0g-ai:9 "git status" C-m
tmux send-keys -t fr0g-ai:9 "git add . && git commit -m 'Feature implementation'" C-m
tmux send-keys -t fr0g-ai:9 "git push origin main" C-m

# Debugging
tmux send-keys -t fr0g-ai:9 "docker-compose logs fr0g-ai-aip" C-m
tmux send-keys -t fr0g-ai:9 "curl -s http://localhost:8080/health | jq" C-m
tmux send-keys -t fr0g-ai:9 "docker ps" C-m

# File operations
tmux send-keys -t fr0g-ai:9 "find . -name '*.go' | grep -v vendor" C-m
tmux send-keys -t fr0g-ai:9 "tail -f logs/fr0g-ai-aip.log" C-m
```

## Advanced Dispatch Patterns

### Cross-Component Coordination

#### Database Migration Campaign
```bash
# 1. Preparation phase
tmux send-keys -t fr0g-ai:1 "Prepare AIP service for database migration" C-m
tmux send-keys -t fr0g-ai:5 "Add database configuration validation" C-m
tmux send-keys -t fr0g-ai:6 "Update Docker Compose with PostgreSQL service" C-m

# 2. Implementation phase
tmux send-keys -t fr0g-ai:1 "Implement database connection pooling" C-m
tmux send-keys -t fr0g-ai:1 "Add transaction support for persona operations" C-m

# 3. Testing phase
tmux send-keys -t fr0g-ai:8 "make test-database-migration" C-m
tmux send-keys -t fr0g-ai:8 "make test-performance" C-m

# 4. Deployment phase
tmux send-keys -t fr0g-ai:6 "Deploy database migration scripts" C-m
tmux send-keys -t fr0g-ai:8 "make health" C-m
```

#### Performance Optimization Campaign
```bash
# System-wide performance optimization
tmux send-keys -t fr0g-ai:1 "Implement caching layer for persona operations" C-m
tmux send-keys -t fr0g-ai:2 "Add connection pooling for external API calls" C-m
tmux send-keys -t fr0g-ai:3 "Optimize cognitive processing algorithms" C-m
tmux send-keys -t fr0g-ai:7 "Optimize service discovery for <5ms response time" C-m

# Performance testing
tmux send-keys -t fr0g-ai:8 "make test-performance" C-m
tmux send-keys -t fr0g-ai:8 "make test-load-1000-users" C-m
```

#### Security Hardening Initiative
```bash
# Security enhancement across all services
tmux send-keys -t fr0g-ai:2 "Implement OAuth2 and JWT authentication" C-m
tmux send-keys -t fr0g-ai:5 "Add security configuration validation" C-m
tmux send-keys -t fr0g-ai:6 "Implement container security scanning" C-m
tmux send-keys -t fr0g-ai:6 "Add network security policies" C-m

# Security validation
tmux send-keys -t fr0g-ai:8 "make validate-production" C-m
tmux send-keys -t fr0g-ai:8 "make test-security" C-m
```

### Agent Communication Protocol

#### Message Types for Coordination
```bash
# Information sharing
tmux send-keys -t fr0g-ai:1 "[INFO] Persona service endpoints ready for integration" C-m
tmux send-keys -t fr0g-ai:2 "[INFO] OpenWebUI integration tests passing" C-m

# Request assistance
tmux send-keys -t fr0g-ai:2 "[REQUEST] Need AIP gRPC endpoint documentation" C-m
tmux send-keys -t fr0g-ai:4 "[REQUEST] Require master-control threat analysis API" C-m

# Task handoff
tmux send-keys -t fr0g-ai:3 "[HANDOFF] Threat analysis logic ready for IO integration" C-m
tmux send-keys -t fr0g-ai:1 "[HANDOFF] Persona CRUD operations complete, ready for Bridge" C-m

# Blocking issues
tmux send-keys -t fr0g-ai:4 "[BLOCKED] Waiting for master-control gRPC interface" C-m
tmux send-keys -t fr0g-ai:2 "[BLOCKED] Need updated AIP service discovery endpoint" C-m

# Completion notifications
tmux send-keys -t fr0g-ai:7 "[COMPLETE] Service registry performance optimization done" C-m
tmux send-keys -t fr0g-ai:1 "[COMPLETE] Database migration successfully implemented" C-m
```

## Development Workflows

### Feature Development Workflow

#### 1. Planning Phase (Project Lead)
```bash
tmux send-keys -t fr0g-ai:0 "Design API endpoints for persona recommendation feature" C-m
tmux send-keys -t fr0g-ai:0 "Update architecture documentation for new feature" C-m
```

#### 2. Implementation Phase (Specialized Agents)
```bash
# Core business logic
tmux send-keys -t fr0g-ai:1 "Implement persona similarity algorithms" C-m
tmux send-keys -t fr0g-ai:1 "Add recommendation engine with machine learning" C-m

# API integration
tmux send-keys -t fr0g-ai:2 "Add recommendation endpoints to Bridge API" C-m
tmux send-keys -t fr0g-ai:2 "Implement caching for recommendation responses" C-m

# Cognitive enhancement
tmux send-keys -t fr0g-ai:3 "Add learning from recommendation feedback" C-m
```

#### 3. Configuration Phase (Config Agent)
```bash
tmux send-keys -t fr0g-ai:5 "Add configuration options for recommendation engine" C-m
tmux send-keys -t fr0g-ai:5 "Implement recommendation algorithm parameters" C-m
```

#### 4. Testing Phase (Build-Test)
```bash
tmux send-keys -t fr0g-ai:8 "make test-recommendation-engine" C-m
tmux send-keys -t fr0g-ai:8 "make test-integration" C-m
tmux send-keys -t fr0g-ai:8 "make test-performance" C-m
```

#### 5. Deployment Phase (DevOps)
```bash
tmux send-keys -t fr0g-ai:6 "Update Docker configuration for recommendation feature" C-m
tmux send-keys -t fr0g-ai:6 "Add monitoring for recommendation endpoints" C-m
```

### Bug Fix Workflow

#### 1. Investigation Phase
```bash
# Check logs and status
tmux send-keys -t fr0g-ai:9 "docker-compose logs | grep ERROR" C-m
tmux send-keys -t fr0g-ai:9 "curl -s http://localhost:8080/health | jq" C-m
tmux send-keys -t fr0g-ai:8 "make health" C-m
```

#### 2. Fix Implementation
```bash
# Identify the responsible agent and dispatch fix
tmux send-keys -t fr0g-ai:1 "Fix validation error in persona service" C-m
tmux send-keys -t fr0g-ai:2 "Resolve OpenWebUI integration timeout issue" C-m
```

#### 3. Testing and Verification
```bash
# Test the fix
tmux send-keys -t fr0g-ai:8 "make test-validation" C-m
tmux send-keys -t fr0g-ai:8 "make test-integration" C-m

# Verify system health
tmux send-keys -t fr0g-ai:8 "make health" C-m
tmux send-keys -t fr0g-ai:9 "curl -s http://localhost:8080/health | jq" C-m
```

### Release Preparation Workflow

#### 1. Code Quality Assurance
```bash
# Code formatting and linting
tmux send-keys -t fr0g-ai:8 "make fmt" C-m
tmux send-keys -t fr0g-ai:8 "make lint" C-m

# Comprehensive testing
tmux send-keys -t fr0g-ai:8 "make test-all" C-m
tmux send-keys -t fr0g-ai:8 "make test-all-integration" C-m
tmux send-keys -t fr0g-ai:8 "make test-performance" C-m
```

#### 2. Security Validation
```bash
tmux send-keys -t fr0g-ai:8 "make validate-production" C-m
tmux send-keys -t fr0g-ai:6 "Run security scanning on all containers" C-m
```

#### 3. Documentation Updates
```bash
tmux send-keys -t fr0g-ai:0 "Update README.md with new features" C-m
tmux send-keys -t fr0g-ai:0 "Update API documentation" C-m
tmux send-keys -t fr0g-ai:5 "Update .env.example with new configuration options" C-m
```

#### 4. Deployment Preparation
```bash
tmux send-keys -t fr0g-ai:6 "Build production Docker images" C-m
tmux send-keys -t fr0g-ai:6 "Test deployment scripts" C-m
tmux send-keys -t fr0g-ai:8 "make docker-build" C-m
```

## Best Practices

### Dispatch System Guidelines

1. **Clear Task Descriptions**: Use specific, actionable task descriptions
2. **Respect Agent Boundaries**: Only dispatch tasks within each agent's domain
3. **Sequential Coordination**: For complex features, dispatch tasks in logical sequence
4. **Monitor Progress**: Check agent windows manually to verify task completion
5. **Use Communication Protocol**: Include message types ([INFO], [REQUEST], etc.) for coordination

### Agent Coordination

1. **Dependency Management**: Ensure prerequisite tasks are completed before dependent tasks
2. **Cross-Component Changes**: Coordinate through Project Lead for changes affecting multiple components
3. **Configuration Changes**: Always involve Config Agent for environment variable changes
4. **Infrastructure Changes**: Always involve DevOps Agent for Docker/deployment changes

### Error Handling

1. **Monitor Agent Windows**: Regularly check agent windows for errors or completion status
2. **Use Shell Window for Debugging**: Window 9 is available for investigation and debugging
3. **Health Checks**: Use Build-Test window (8) for system health verification
4. **Rollback Procedures**: Have rollback commands ready for failed deployments

### Performance Optimization

1. **Batch Related Tasks**: Group related tasks to the same agent for efficiency
2. **Parallel Execution**: Dispatch independent tasks to different agents simultaneously
3. **Resource Management**: Monitor system resources when running multiple agents
4. **Task Prioritization**: Handle critical bugs and security issues first

## Troubleshooting

### Common Issues

#### Agent Not Responding
```bash
# Check if agent window is active
tmux list-windows -t fr0g-ai

# Switch to agent window manually
tmux select-window -t fr0g-ai:1

# Restart agent if needed
tmux send-keys -t fr0g-ai:1 C-c
tmux send-keys -t fr0g-ai:1 "aider --message-file .aider-prompts/aip-agent.md TODO.md" C-m
```

#### Command Not Executing
```bash
# Verify tmux session exists
tmux list-sessions

# Check window status
tmux list-windows -t fr0g-ai

# Send command with explicit session
tmux send-keys -t fr0g-ai:1 "echo 'test command'" C-m
```

#### Agent Working on Wrong Component
```bash
# Check current directory in agent window
tmux send-keys -t fr0g-ai:1 "pwd" C-m

# Navigate to correct directory
tmux send-keys -t fr0g-ai:1 "cd fr0g-ai-aip" C-m
```

### Recovery Procedures

#### Restart Development Environment
```bash
# Kill existing session
tmux kill-session -t fr0g-ai

# Restart development environment
./start-fr0g-ai-dev.sh
```

#### Reset Specific Agent
```bash
# Kill specific window
tmux kill-window -t fr0g-ai:1

# Recreate window
tmux new-window -t fr0g-ai -n "AIP" -c "fr0g-ai-aip"
tmux send-keys -t fr0g-ai:1 "aider --message-file .aider-prompts/aip-agent.md TODO.md" C-m
```

This comprehensive dispatch system enables sophisticated coordination across the entire fr0g.ai development ecosystem while maintaining clear separation of concerns and specialized expertise.
