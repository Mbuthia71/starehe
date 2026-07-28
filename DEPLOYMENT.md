# VPS Deployment Guide

## Prerequisites
- SSH access to your VPS
- Docker and Docker Compose installed on VPS
- Project files copied to VPS

## Deployment Steps

### 1. Copy Project Files to VPS
```bash
# From your local machine
scp -r . user@your-vps-ip:/path/to/starehian-society-platform
```

### 2. SSH into VPS
```bash
ssh user@your-vps-ip
cd /path/to/starehian-society-platform
```

### 3. Create .env File
```bash
cp .env.example .env
nano .env
```

Update the following environment variables with your actual values:
- POSTGRES_PASSWORD
- JWT_SECRET
- REFRESH_TOKEN_SECRET
- AFRICAS_TALKING_API_KEY
- AFRICAS_TALKING_USERNAME
- CLOUDFLARE_R2_ACCESS_KEY
- CLOUDFLARE_R2_SECRET_KEY
- CLOUDFLARE_R2_BUCKET
- CLOUDFLARE_R2_ENDPOINT
- FCM_SERVER_KEY
- CENTRIFUGO_SECRET
- CENTRIFUGO_API_KEY
- GRAFANA_ADMIN_PASSWORD

### 4. Stop Existing Containers (if any)
```bash
docker-compose down
```

### 5. Build and Start Services
```bash
# Build new images
docker-compose build api frontend

# Start all services
docker-compose up -d

# Or start specific services
docker-compose up -d postgres redis centrifugo api frontend
```

### 6. Check Service Status
```bash
docker-compose ps
docker-compose logs -f api
docker-compose logs -f frontend
```

### 7. Verify Deployment
```bash
# Check API health
curl http://localhost:3000/health

# Check frontend
curl http://localhost:3001
```

## Useful Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f frontend
docker-compose logs -f postgres
```

### Restart Services
```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart api
docker-compose restart frontend
```

### Update and Redeploy
```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose build api frontend
docker-compose up -d api frontend
```

### Database Migrations
The migrations are automatically run on PostgreSQL startup via the volume mount in docker-compose.yml.

### Backup Database
```bash
docker-compose run backup
```

## Service Ports
- API: 3000
- Frontend: 3001
- PostgreSQL: 5432
- Redis: 6379
- Centrifugo: 8000
- Grafana: 3000 (mapped to 3001 in compose, may conflict - adjust if needed)
- Loki: 3100

## Troubleshooting

### Container won't start
```bash
docker-compose logs [service-name]
```

### Database connection issues
```bash
docker-compose exec postgres psql -U starehian_user -d starehian_db
```

### Redis connection issues
```bash
docker-compose exec redis redis-cli ping
```

### Clear everything and start fresh
```bash
docker-compose down -v
docker-compose up -d
```

## Production Considerations
- Use strong passwords in .env
- Configure firewall rules
- Set up SSL/TLS with nginx reverse proxy
- Monitor logs regularly
- Set up automated backups
