# API Documentation

## Overview

fr0g.ai exposes both HTTP REST APIs and gRPC APIs for different use cases:

- **HTTP APIs**: External integrations, web interfaces, simple clients
- **gRPC APIs**: High-performance inter-service communication, advanced clients

## Service APIs

### fr0g-ai-aip (Core AI Service)

**Base URLs:**
- HTTP: `http://localhost:8080`
- gRPC: `localhost:9090`

**Key Endpoints:**
- `GET /health` - Health check
- `POST /api/process` - Process AI requests
- `GET /api/status` - Service status
- `GET /metrics` - Prometheus metrics

### fr0g-ai-bridge (Integration Bridge)

**Base URLs:**
- HTTP: `http://localhost:8081`
- gRPC: `localhost:9091`

**Key Endpoints:**
- `GET /health` - Health check
- `POST /api/bridge` - Bridge requests to external systems
- `GET /api/integrations` - List available integrations
- `GET /metrics` - Prometheus metrics

## Authentication

### API Keys
```bash
# Include API key in headers
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/process
```

### JWT Tokens
```bash
# Include JWT token in Authorization header
curl -H "Authorization: Bearer your-jwt-token" http://localhost:8080/api/process
```

## Error Handling

### HTTP Status Codes
- `200` - Success
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `429` - Rate Limited
- `500` - Internal Server Error
- `503` - Service Unavailable

### Error Response Format
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Request validation failed",
    "details": {
      "field": "email",
      "reason": "invalid format"
    }
  }
}
```

## Rate Limiting

All APIs implement rate limiting:
- **Default**: 100 requests per minute per client
- **Headers**: Rate limit info included in response headers
- **Exceeded**: Returns 429 status with retry information

```bash
# Rate limit headers
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## Examples

### Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2023-12-01T12:00:00Z",
  "version": "1.0.0",
  "uptime": "24h30m15s"
}
```

### Process AI Request
```bash
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "type": "email_analysis",
    "data": {
      "sender": "unknown@example.com",
      "subject": "Urgent: Account Verification Required",
      "body": "Click here to verify your account..."
    }
  }'
```

Response:
```json
{
  "id": "req_123456789",
  "status": "completed",
  "result": {
    "threat_level": "high",
    "confidence": 0.95,
    "threats": ["phishing", "social_engineering"],
    "recommendation": "block"
  },
  "processing_time_ms": 150
}
```

## gRPC APIs

### Service Definitions

Services are defined in Protocol Buffer files:
- `fr0g-ai-aip/proto/aip.proto`
- `fr0g-ai-bridge/proto/bridge.proto`

### Client Examples

#### Go Client
```go
package main

import (
    "context"
    "google.golang.org/grpc"
    pb "github.com/fr0g-vibe/fr0g-ai-aip/proto"
)

func main() {
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pb.NewAIPServiceClient(conn)
    
    req := &pb.ProcessRequest{
        Type: "email_analysis",
        Data: map[string]string{
            "sender": "unknown@example.com",
            "subject": "Urgent: Account Verification Required",
        },
    }
    
    resp, err := client.Process(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Threat Level: %s\n", resp.ThreatLevel)
}
```

#### Python Client
```python
import grpc
import aip_pb2
import aip_pb2_grpc

def main():
    with grpc.insecure_channel('localhost:9090') as channel:
        stub = aip_pb2_grpc.AIPServiceStub(channel)
        
        request = aip_pb2.ProcessRequest(
            type='email_analysis',
            data={
                'sender': 'unknown@example.com',
                'subject': 'Urgent: Account Verification Required'
            }
        )
        
        response = stub.Process(request)
        print(f"Threat Level: {response.threat_level}")

if __name__ == '__main__':
    main()
```

## WebSocket APIs (Future)

Real-time communication for streaming data:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function() {
    ws.send(JSON.stringify({
        type: 'subscribe',
        channel: 'threat_alerts'
    }));
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Threat Alert:', data);
};
```

## SDK Libraries (Future)

Official SDKs will be available for:
- Go
- Python
- JavaScript/Node.js
- Java
- C#

## Testing APIs

### Using curl
```bash
# Health check
curl http://localhost:8080/health

# Process request with authentication
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test-key" \
  -d '{"type": "test", "data": {}}'
```

### Using grpcurl
```bash
# List services
grpcurl -plaintext localhost:9090 list

# Call method
grpcurl -plaintext -d '{"type": "test"}' localhost:9090 aip.AIPService/Process
```

### Using Postman

Import the provided Postman collection:
- `docs/api/fr0g-ai.postman_collection.json`

## API Versioning

APIs follow semantic versioning:
- **v1**: Current stable version
- **v2**: Future version with breaking changes

Version is specified in URL path:
- `http://localhost:8080/api/v1/process`
- `http://localhost:8080/api/v2/process`

## OpenAPI Specification

Full OpenAPI 3.0 specifications available:
- [AIP Service OpenAPI](aip-openapi.yaml)
- [Bridge Service OpenAPI](bridge-openapi.yaml)

## Support

For API support:
- Check service logs for error details
- Review this documentation
- Submit issues on GitHub
- Contact support team
