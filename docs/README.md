# fr0g.ai Documentation

## Platform Overview

fr0g.ai is a comprehensive microservices-based AI security platform consisting of five main components working together to provide intelligent threat detection, response automation, and cognitive orchestration.

### [Service Registry Component](./service-registry-component.md)
Central service discovery and health monitoring system.
- **Engineers**: Focus on service discovery, load balancing, and health monitoring
- **Key Files**: Registry server, service discovery APIs, health monitoring
- **Status**: OPERATIONAL - 9,553+ ops/sec performance

### [AI Persona (AIP) Component](./ai-persona-component.md)
Core AI processing engine for persona and identity management.
- **Engineers**: Focus on AI integration, persona behavior, and rich attributes processing
- **Key Files**: Persona processors, identity management, storage systems
- **Status**: OPERATIONAL - 293 personas, 8 attribute processors

### [Bridge Component](./bridge-component.md) 
Integration layer providing HTTP REST and gRPC interfaces for external clients.
- **Engineers**: Focus on API design, protocol handling, and external integrations
- **Key Files**: HTTP handlers, gRPC services, middleware, OpenWebUI integration
- **Status**: OPERATIONAL - OpenWebUI integration verified

### [Master Control Program (MCP)](./master-control-component.md)
Central intelligence and orchestration engine with conscious AI capabilities.
- **Engineers**: Focus on cognitive architecture, intelligent orchestration, and emergent capabilities
- **Key Files**: Cognitive engines, learning systems, pattern recognition, consciousness
- **Status**: OPERATIONAL - Conscious AI with 0.154 learning rate

### [Input/Output (I/O) Component](./io-component.md)
Comprehensive I/O processing for all threat vectors and external communications.
- **Engineers**: Focus on threat vector processing, external API integration, and response automation
- **Key Files**: Input processors, output managers, threat detection, message queuing
- **Status**: OPERATIONAL - All 5 input processors working

## Inter-Component Communication

```
External Clients
       ↓
   Bridge Component ←→ Service Registry ←→ Master Control Program (MCP)
       ↓                      ↓                      ↓
   AI Persona Component ←─────┼──────────────────────┘
       ↑                      ↓
       └── I/O Component ←────┘
              ↓
       Threat Vectors
    (SMS, Voice, Email, IRC, Discord)
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
