# fr0g.ai Port Assignment Documentation

## Service Port Assignments

### Production Port Mapping
All services are configured with specific ports to avoid conflicts:

| Service | HTTP Port | gRPC Port | Docker Mapping | Status |
|---------|-----------|-----------|----------------|---------|
| fr0g-ai-registry | 8500 | N/A | 8500:8500 | ✅ Operational |
| fr0g-ai-aip | 8080 | 9090 | 8080:8080, 9090:9090 | ✅ Operational |
| fr0g-ai-bridge | 8082 | 9092 | 8082:8082, 9092:9092 | ✅ Operational |
| fr0g-ai-master-control | 8081 | N/A | 8081:8081 | ❌ Needs HTTP server |
| fr0g-ai-io | 8083 | 9093 | 8083:8083, 9093:9093 | ❌ Port mismatch fixed |

### Port Conflict Resolution
- **No conflicts**: All services use unique ports
- **Sequential assignment**: Ports assigned in logical sequence
- **Docker mapping**: 1:1 mapping between host and container ports

### Health Check Endpoints
All services expose health endpoints on their HTTP ports:

```bash
# Service Registry
curl http://localhost:8500/health

# AIP Service  
curl http://localhost:8080/health

# Bridge Service
curl http://localhost:8082/health

# Master Control (after fix)
curl http://localhost:8081/health

# I/O Service (after fix)
curl http://localhost:8083/health
```

### Configuration Sources
Port assignments are defined in:
1. `docker-compose.yml` - Container port mappings
2. Service environment variables - Internal binding
3. `Makefile` health-quick target - Health check verification

### Service Communication
Inter-service communication uses Docker network names:
- `fr0g-ai-aip:9090` (gRPC)
- `fr0g-ai-bridge:9092` (gRPC)
- `fr0g-ai-master-control:8081` (HTTP)
- `fr0g-ai-io:9093` (gRPC)
- `service-registry:8500` (HTTP)

### Issues Resolved
1. **fr0g-ai-io**: Fixed port binding from 8080/9090 to 8083/9093
2. **fr0g-ai-master-control**: Added HTTP server for health checks
3. **Health checks**: Verified all ports match docker-compose configuration

### Verification Commands
```bash
# Check all service health
make health-quick

# Check Docker container status
docker-compose ps

# Check individual service logs
docker-compose logs fr0g-ai-master-control
docker-compose logs fr0g-ai-io
```
