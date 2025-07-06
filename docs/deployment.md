# Deployment Guide

## Overview

fr0g.ai is designed for flexible deployment across different environments, from local development to production Kubernetes clusters.

## Deployment Options

### 1. Docker Compose (Recommended for Development/Testing)

#### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- 2GB+ RAM
- 10GB+ disk space

#### Quick Deployment
```bash
# Clone and setup
git clone --recursive https://github.com/fr0g-vibe/fr0g-ai.git
cd fr0g-ai

# Configure environment
cp .env.example .env
# Edit .env with your configuration

# Deploy
make setup
docker-compose up -d

# Verify deployment
make health
```

#### Production Docker Compose
```bash
# Use production configuration
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 2. Kubernetes (Production)

#### Prerequisites
- Kubernetes 1.20+
- kubectl configured
- Helm 3.0+ (optional)

#### Manual Deployment
```bash
# Create namespace
kubectl create namespace fr0g-ai

# Deploy services
kubectl apply -f k8s/

# Check status
kubectl get pods -n fr0g-ai
```

#### Helm Deployment (Future)
```bash
# Add fr0g.ai Helm repository
helm repo add fr0g-ai https://charts.fr0g.ai
helm repo update

# Install
helm install fr0g-ai fr0g-ai/fr0g-ai -n fr0g-ai --create-namespace
```

### 3. Cloud Platforms

#### AWS ECS
- Use provided ECS task definitions
- Configure Application Load Balancer
- Set up CloudWatch logging

#### Google Cloud Run
- Deploy containerized services
- Configure Cloud Load Balancing
- Use Cloud Logging

#### Azure Container Instances
- Deploy using ARM templates
- Configure Azure Load Balancer
- Use Azure Monitor

## Configuration

### Environment Variables

#### Core Configuration
```bash
# Service Configuration
FR0G_AIP_STORAGE_TYPE=file|database
FR0G_AIP_DATA_DIR=/app/data
FR0G_BRIDGE_HTTP_PORT=8080
FR0G_BRIDGE_GRPC_PORT=9090

# Logging
LOG_LEVEL=info|debug|warn|error
LOG_FORMAT=json|text

# Security
JWT_SECRET=your-jwt-secret
WEBUI_SECRET_KEY=your-webui-secret
```

#### External Integrations
```bash
# OpenWebUI Integration
OPENWEBUI_API_KEY=your-api-key
OPENWEBUI_BASE_URL=https://your-openwebui-instance.com

# Database (if using database storage)
DATABASE_URL=postgres://user:pass@host:5432/fr0g_ai

# Redis (for caching/queues)
REDIS_URL=redis://host:6379
```

### Storage Configuration

#### File Storage (Default)
```bash
FR0G_AIP_STORAGE_TYPE=file
FR0G_AIP_DATA_DIR=/app/data
```

#### Database Storage
```bash
FR0G_AIP_STORAGE_TYPE=database
DATABASE_URL=postgres://user:pass@host:5432/fr0g_ai
```

### Networking

#### Port Configuration
- **fr0g-ai-registry**: HTTP 8500
- **fr0g-ai-aip**: HTTP 8080, gRPC 9090
- **fr0g-ai-bridge**: HTTP 8082, gRPC 9091
- **fr0g-ai-master-control**: HTTP 8081
- **fr0g-ai-io**: HTTP 8083, gRPC 9092

#### Load Balancer Configuration
```nginx
# Nginx example
upstream fr0g-ai-registry {
    server fr0g-ai-registry:8500;
}

upstream fr0g-ai-aip {
    server fr0g-ai-aip:8080;
}

upstream fr0g-ai-bridge {
    server fr0g-ai-bridge:8082;
}

upstream fr0g-ai-master-control {
    server fr0g-ai-master-control:8081;
}

upstream fr0g-ai-io {
    server fr0g-ai-io:8083;
}

server {
    listen 80;
    server_name api.fr0g.ai;
    
    location /registry/ {
        proxy_pass http://fr0g-ai-registry/;
    }
    
    location /aip/ {
        proxy_pass http://fr0g-ai-aip/;
    }
    
    location /bridge/ {
        proxy_pass http://fr0g-ai-bridge/;
    }
    
    location /control/ {
        proxy_pass http://fr0g-ai-master-control/;
    }
    
    location /io/ {
        proxy_pass http://fr0g-ai-io/;
    }
}
```

## Security

### TLS/SSL Configuration

#### Docker Compose with Traefik
```yaml
# docker-compose.prod.yml
services:
  traefik:
    image: traefik:v2.9
    command:
      - --api.dashboard=true
      - --entrypoints.web.address=:80
      - --entrypoints.websecure.address=:443
      - --certificatesresolvers.letsencrypt.acme.email=admin@fr0g.ai
      - --certificatesresolvers.letsencrypt.acme.storage=/acme.json
      - --certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./acme.json:/acme.json
    labels:
      - traefik.enable=true
      - traefik.http.routers.api.rule=Host(`traefik.fr0g.ai`)
      - traefik.http.routers.api.tls.certresolver=letsencrypt

  fr0g-ai-aip:
    labels:
      - traefik.enable=true
      - traefik.http.routers.aip.rule=Host(`aip.fr0g.ai`)
      - traefik.http.routers.aip.tls.certresolver=letsencrypt
```

### Secrets Management

#### Docker Secrets
```yaml
# docker-compose.yml
secrets:
  jwt_secret:
    file: ./secrets/jwt_secret.txt
  api_key:
    file: ./secrets/api_key.txt

services:
  fr0g-ai-bridge:
    secrets:
      - jwt_secret
      - api_key
```

#### Kubernetes Secrets
```bash
# Create secrets
kubectl create secret generic fr0g-ai-secrets \
  --from-literal=jwt-secret=your-jwt-secret \
  --from-literal=api-key=your-api-key \
  -n fr0g-ai
```

## Monitoring

### Health Checks

All services expose health endpoints:
```bash
# Health check endpoints
curl http://localhost:8080/health  # AIP service
curl http://localhost:8081/health  # Bridge service
```

### Metrics

#### Prometheus Configuration
```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'fr0g-ai-aip'
    static_configs:
      - targets: ['fr0g-ai-aip:8080']
    metrics_path: /metrics

  - job_name: 'fr0g-ai-bridge'
    static_configs:
      - targets: ['fr0g-ai-bridge:8081']
    metrics_path: /metrics
```

### Logging

#### Centralized Logging with ELK Stack
```yaml
# docker-compose.monitoring.yml
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.5.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false

  logstash:
    image: docker.elastic.co/logstash/logstash:8.5.0
    volumes:
      - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf

  kibana:
    image: docker.elastic.co/kibana/kibana:8.5.0
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
```

## Scaling

### Horizontal Scaling

#### Docker Compose
```bash
# Scale services
docker-compose up -d --scale fr0g-ai-aip=3 --scale fr0g-ai-bridge=2
```

#### Kubernetes
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fr0g-ai-aip
spec:
  replicas: 3
  selector:
    matchLabels:
      app: fr0g-ai-aip
  template:
    metadata:
      labels:
        app: fr0g-ai-aip
    spec:
      containers:
      - name: fr0g-ai-aip
        image: fr0g-ai-aip:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090
```

### Auto-scaling

#### Kubernetes HPA
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: fr0g-ai-aip-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: fr0g-ai-aip
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

## Backup and Recovery

### Data Backup

#### File Storage Backup
```bash
# Backup data directory
tar -czf fr0g-ai-backup-$(date +%Y%m%d).tar.gz data/

# Restore
tar -xzf fr0g-ai-backup-20231201.tar.gz
```

#### Database Backup
```bash
# PostgreSQL backup
pg_dump $DATABASE_URL > fr0g-ai-backup-$(date +%Y%m%d).sql

# Restore
psql $DATABASE_URL < fr0g-ai-backup-20231201.sql
```

### Disaster Recovery

1. **Service Recovery**: Use health checks and auto-restart policies
2. **Data Recovery**: Regular automated backups
3. **Configuration Recovery**: Store configuration in version control
4. **Infrastructure Recovery**: Infrastructure as Code (Terraform/CloudFormation)

## Performance Tuning

### Resource Limits

#### Docker Compose
```yaml
services:
  fr0g-ai-aip:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '1.0'
          memory: 1G
```

#### Kubernetes
```yaml
resources:
  limits:
    cpu: 2000m
    memory: 2Gi
  requests:
    cpu: 1000m
    memory: 1Gi
```

### Database Optimization

```sql
-- PostgreSQL optimization
CREATE INDEX CONCURRENTLY idx_data_timestamp ON data_table(timestamp);
CREATE INDEX CONCURRENTLY idx_data_type ON data_table(type);
```

## Troubleshooting

### Common Issues

1. **Service Won't Start**
   - Check logs: `docker-compose logs service-name`
   - Verify configuration
   - Check port availability

2. **High Memory Usage**
   - Monitor with `docker stats`
   - Adjust resource limits
   - Check for memory leaks

3. **Slow Performance**
   - Check CPU/memory usage
   - Review database queries
   - Optimize network configuration

### Debug Commands

```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs -f

# Execute commands in container
docker-compose exec fr0g-ai-aip /bin/sh

# Check resource usage
docker stats
```

## Maintenance

### Updates

```bash
# Update submodules
make update-submodules

# Rebuild and restart
docker-compose build
docker-compose up -d

# Verify health
make health
```

### Cleanup

```bash
# Clean old containers and images
make docker-clean

# Clean system
docker system prune -a
```
