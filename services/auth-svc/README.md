# Auth Service

Authentication service for Bitzap platform.

## Setup

### 1. Environment Variables

Copy the example environment file and configure your variables:

```bash
cp env.example .env
```

Edit `.env` file with your actual values:

```bash
# Email Configuration
MAILJET_API_KEY=your-actual-mailjet-api-key
MAILJET_SECRET_KEY=your-actual-mailjet-secret-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your App Name
APP_URL=http://localhost:8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=auth_service
DB_SSL_MODE=disable

# Redis Configuration
REDIS_ADDRESS=127.0.0.1:6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# Server Configuration
SERVER_PORT=8080
LOG_LEVEL=info
```

### 2. Database Setup

Make sure PostgreSQL is running and create the database:

```sql
CREATE DATABASE auth_service;
```

### 3. Redis Setup

Make sure Redis is running on the configured address.

### 4. Run the Application

```bash
go mod tidy
go run cmd/main.go
```

## API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/
```

## Security Notes

- Never commit `.env` file to version control
- Use strong passwords for database and Redis
- Keep your Mailjet API keys secure
- Use HTTPS in production
- Regularly rotate your API keys

## Environment Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `MAILJET_API_KEY` | Mailjet API key for sending emails | - |
| `MAILJET_SECRET_KEY` | Mailjet secret key | - |
| `FROM_EMAIL` | Default sender email | noreply@bitzap.com |
| `FROM_NAME` | Default sender name | Bitzap Auth Service |
| `APP_URL` | Application base URL | http://localhost:8080 |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database username | admin |
| `DB_PASSWORD` | Database password | admin123 |
| `DB_NAME` | Database name | auth_service |
| `REDIS_ADDRESS` | Redis address | 127.0.0.1:6379 |
| `REDIS_PASSWORD` | Redis password | redispass |
| `SERVER_PORT` | Server port | 8080 |
| `LOG_LEVEL` | Log level | info | 