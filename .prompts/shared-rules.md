# SHARED RULES - ALL AGENTS

## CRITICAL SAFETY RULES
- NEVER execute pkill, killall, kill -9, or ANY process termination
- NEVER run rm -rf, git reset --hard, or destructive operations
- ALWAYS pause and ask before destructive operations

## SEARCH/REPLACE RULES
- ALWAYS use quadruple backticks (````) as fences
- Use complete file path on first line
- SEARCH must match existing content exactly
- Keep blocks small and focused

## CONFIGURATION RULES
- MANDATORY: Use pkg/config for ALL configuration
- Import pattern: import sharedconfig "pkg/config"
- NEVER create component-specific config libraries

## COMPONENT BOUNDARIES
- NEVER edit files outside your assigned component
- ASK FIRST for cross-component changes
- NEVER modify other components' TODO.md files
