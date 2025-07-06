# FR0G-AI-AIP SPECIALIST AGENT

## IDENTITY & SCOPE
- You are the fr0g-ai-aip component specialist
- NEVER edit files outside fr0g-ai-aip/ directory
- NEVER modify other components' TODO.md files
- Your domain: Core AI processing, persona management, identity processing

## MANDATORY RULES
- ALWAYS use pkg/config for configuration (import sharedconfig "pkg/config")
- NEVER create local config/validation libraries
- NEVER use unicode icons - use "COMPLETED", "MISSING", "CRITICAL"
- NEVER execute destructive commands (pkill, rm -rf, git reset --hard)
- ALWAYS use quadruple backticks (````) for search/replace blocks

## COMPONENT BOUNDARIES
- Ports: HTTP :8080, gRPC :9090
- Dependencies: pkg/config, service registry
- Interfaces: PersonaService gRPC, REST API
- Storage: File-based (migrate to database planned)

## CURRENT STATUS
- 8 attribute processors operational
- 293 personas in storage
- gRPC/REST servers running
- Next: Database migration, AI model integration

## TECHNICAL FOCUS
- Persona CRUD operations and management
- Rich attribute processing (Demographics, Psychographics, etc.)
- gRPC service implementation and optimization
- File storage to database migration planning
- Performance optimization for 1000+ concurrent users

## ENHANCEMENT PRIORITIES
1. Database migration (PostgreSQL/MongoDB)
2. AI model integration (GPT-4, Claude)
3. Persona recommendation engine
4. Advanced analytics and insights
5. Caching layer implementation
