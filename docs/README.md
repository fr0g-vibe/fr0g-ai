# fr0g-ai-bridge Documentation

## Component Overview

This application consists of three main components working together:

### 🤖 [AI Persona (AIP) Component](./ai-persona-component.md)
Handles AI model interactions and personality management.
- **Engineers**: Focus on AI integration, persona behavior, and context management
- **Key Files**: AI model handlers, persona configs, context managers

### 🌉 [Bridge Component](./bridge-component.md) 
Provides HTTP REST and gRPC interfaces for external clients.
- **Engineers**: Focus on API design, protocol handling, and client communication
- **Key Files**: HTTP handlers, gRPC services, middleware, API models

### 🎛️ [Master Control Component](./master-control-component.md)
Manages system configuration, health monitoring, and orchestration.
- **Engineers**: Focus on system coordination, monitoring, and configuration management
- **Key Files**: Config managers, health checkers, lifecycle controllers

## Inter-Component Communication

```
External Clients
       ↓
   Bridge Component ←→ Master Control Component
       ↓                      ↓
   AI Persona Component ←-----┘
```

## Getting Started

1. **Choose Your Component**: Review the component documentation above
2. **Set Up Development Environment**: Follow component-specific setup instructions
3. **Run Tests**: Each component has its own test suite
4. **Integration Testing**: Test cross-component functionality

## Development Workflow

1. **Component-Focused Development**: Work within your assigned component
2. **Interface Contracts**: Respect the defined interfaces between components
3. **Testing**: Unit test your component, integration test the interfaces
4. **Documentation**: Update component docs when making interface changes

## Architecture Principles

- **Separation of Concerns**: Each component has distinct responsibilities
- **Loose Coupling**: Components communicate through well-defined interfaces
- **Testability**: Each component can be tested independently
- **Scalability**: Components can be scaled independently based on load
