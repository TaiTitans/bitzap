# Kong Gateway Setup for BitZap URL Shortener

Kong Gateway is configured as the API Gateway for the BitZap URL Shortener project with the following features:

- **Authentication**: JWT-based authentication
- **Rate Limiting**: Request limits per minute/hour
- **CORS**: Cross-Origin Resource Sharing
- **Logging**: Request/Response logging
- **Load Balancing**: Distributes traffic among services

## 🚀 Quick Start

### 1. Start Services
```bash
docker-compose up -d
```

### 2. Setup Kong Gateway
```bash
chmod +x infra/kong/setup-kong.sh
./infra/kong/setup-kong.sh
```

### 3. Generate JWT Token
```bash
python3 infra/kong/generate-jwt.py
```

## 📋 API Endpoints

### Public Endpoints (No JWT Required)
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register
- `GET /r/{short_code}` - Redirect URL

### Protected Endpoints (JWT Required)
- `POST /api/v1/shorten` - Create short URL
- `GET /api/v1/urls` - List URLs
- `PUT /api/v1/urls` - Update URL
- `DELETE /api/v1/urls` - Delete URL
- `GET /api/v1/analytics` - Analytics dashboard
- `GET /api/v1/billing` - Billing information

## 🔑 Authentication

### Why JWT?
Kong Gateway uses JWT (JSON Web Token) to:
- **Authenticate users**: Ensure only valid users can access protected APIs
- **Stateless**: No need to store sessions, token contains all necessary info
- **Secure**: Token is signed with a secret key, cannot be forged
- **Scalable**: Different services can verify tokens independently

### JWT Consumers
- **App Consumer**: `bitzap-key` / `bitzap-secret-2025`
- **Mobile Consumer**: `mobile-key` / `mobile-secret-2025`

### Generate JWT Token
```bash
# Using Go
cd infra/kong && go run generate-jwt.go
```

### Using JWT Token
```bash
# Generate token
cd infra/kong && go run generate-jwt.go

# Copy the token and use it in API calls
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com"}' \
     http://localhost:8000/api/v1/shorten
```

### JWT Token Structure
```json
{
  "iss": "bitzap-key",           // Issuer (Kong key)
  "sub": "testuser",             // Subject (user ID)
  "iat": 1640995200,              // Issued at
  "exp": 1641081600,              // Expires at
  "user_id": "testuser",         // Custom claim
  "username": "testuser"         // Custom claim
}
```

## 🛠️ Management

### Kong Admin API
- **URL**: http://localhost:8001
- **Health Check**: `curl http://localhost:8001/status`

### Konga UI (Web Interface)
- **URL**: http://localhost:1337
- **Setup**: Connect to Kong Admin API at `http://kong:8001`

## 📊 Monitoring

### View Logs
```bash
# Kong logs
docker logs kong

# Service logs
docker logs auth-svc
docker logs shortener-svc
```

### Check Kong Status
```bash
curl http://localhost:8001/status
```

### List Services
```bash
curl http://localhost:8001/services
```

### List Routes
```bash
curl http://localhost:8001/routes
```

## 🔧 Configuration

### Kong Configuration File
The `kong.yml` file contains all service, route, consumer, and plugin configurations.

### Environment Variables
```yaml
# Kong Database
KONG_DATABASE: postgres
KONG_PG_HOST: kong-database
KONG_PG_PASSWORD: kong
KONG_PG_USER: kong
KONG_PG_DATABASE: kong

# Logging
KONG_PROXY_ACCESS_LOG: /dev/stdout
KONG_ADMIN_ACCESS_LOG: /dev/stdout
KONG_PROXY_ERROR_LOG: /dev/stderr
KONG_ADMIN_ERROR_LOG: /dev/stderr
```

## 🚨 Troubleshooting

### Kong fails to start
```bash
# Check database connection
docker logs kong-database

# Check Kong migrations
docker logs kong-migrations

# Restart Kong
docker-compose restart kong
```

### JWT Authentication Error
```bash
# Check consumers
curl http://localhost:8001/consumers

# Check JWT credentials
curl http://localhost:8001/consumers/bitzap-app/jwt

# Regenerate token
python3 infra/kong/generate-jwt.py
```

### Service not accessible
```bash
# Check service health
curl http://localhost:8001/services/auth-service

# Check routes
curl http://localhost:8001/services/auth-service/routes

# Test direct service
curl http://localhost:8081/health  # auth-svc
curl http://localhost:8082/health  # shortener-svc
```

## 📈 Performance Tuning

### Rate Limiting
- **Minute**: 100 requests
- **Hour**: 1000 requests  
- **Day**: 10000 requests

### Request Size Limiting
- **Max Payload**: 10MB

### Caching (Optional)
You can enable the Redis caching plugin:
```bash
curl -X POST http://localhost:8001/plugins \
  --data name=proxy-cache \
  --data config.response_code=200 \
  --data config.request_method=GET \
  --data config.content_type=application/json \
  --data config.cache_ttl=300
```

## 🔐 Security Best Practices

1. **JWT Secret**: Use strong secrets and rotate them regularly
2. **HTTPS**: Enable SSL in production
3. **Rate Limiting**: Adjust limits according to your traffic
4. **IP Restriction**: Restrict admin API access
5. **Logging**: Monitor and alert on suspicious activities

## 📚 Useful Commands

```bash
# Reload Kong configuration
curl -X POST http://localhost:8001/config \
  --form config=@infra/kong/kong.yml

# Add new consumer
curl -X POST http://localhost:8001/consumers \
  --data username=new-user

# Enable plugin on specific route
curl -X POST http://localhost:8001/routes/ROUTE_ID/plugins \
  --data name=rate-limiting \
  --data config.minute=50

# Check plugin status
curl http://localhost:8001/plugins
```