# fr0g.ai Demos

This directory contains standalone demonstrations of various fr0g.ai system components.

## Available Demos

### Webhook Demo (`webhook-demo/`)
Demonstrates the webhook input system functionality of the Master Control Program.

**Features:**
- Tests webhook endpoints for different platforms (Discord, etc.)
- Shows health and status monitoring
- Demonstrates AI community review process
- Interactive testing of webhook processing

**Usage:**
```bash
cd demos/webhook-demo
go run main.go
```

## Running Demos

Each demo is self-contained and can be run independently. Make sure you have the necessary dependencies installed and any required services running.

## Adding New Demos

When adding new demos:
1. Create a new directory under `demos/`
2. Include a `main.go` file with the demo logic
3. Add documentation to this README
4. Ensure the demo is self-contained and doesn't interfere with production systems
