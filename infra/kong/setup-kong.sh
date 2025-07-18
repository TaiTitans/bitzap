#!/bin/bash

# Kong Gateway Setup Script for BitZap URL Shortener
# This script configures Kong with services, routes, and authentication

KONG_ADMIN_URL="http://localhost:8001"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Setting up Kong Gateway for BitZap..."

# Wait for Kong to be ready
echo "Waiting for Kong to be ready..."
until curl -f $KONG_ADMIN_URL/status > /dev/null 2>&1; do
    echo "Waiting for Kong..."
    sleep 5
done

echo "Kong is ready!"

# Method 1: Load configuration from YAML file (recommended)
if command -v deck &> /dev/null; then
    echo "Loading configuration from kong.yml using deck..."
    deck sync --kong-addr $KONG_ADMIN_URL --state $SCRIPT_DIR/kong.yml
    echo "Configuration loaded from YAML!"
else
    echo "deck not found, using manual API calls..."
    
    # Method 2: Manual API configuration
    echo "Creating services manually..."

    # Auth Service
    curl -i -X POST $KONG_ADMIN_URL/services \
      --data name=auth-service \
      --data url=http://auth-svc:8080

    # Shortener Service  
    curl -i -X POST $KONG_ADMIN_URL/services \
      --data name=shortener-service \
      --data url=http://shortener-svc:8080

    # Redirect Service
    curl -i -X POST $KONG_ADMIN_URL/services \
      --data name=redirect-service \
      --data url=http://redirect-svc:8080

    # Analytics Service
    curl -i -X POST $KONG_ADMIN_URL/services \
      --data name=analytics-service \
      --data url=http://analytics-svc:8080

    # Billing Service
    curl -i -X POST $KONG_ADMIN_URL/services \
      --data name=billing-service \
      --data url=http://billing-svc:8080

    echo "Services created!"

    # Create Routes
    echo "Creating routes..."

    # Auth routes (public)
    curl -i -X POST $KONG_ADMIN_URL/services/auth-service/routes \
      --data 'paths[]=/api/v1/auth/login' \
      --data 'methods[]=POST' \
      --data name=auth-login

    curl -i -X POST $KONG_ADMIN_URL/services/auth-service/routes \
      --data 'paths[]=/api/v1/auth/register' \
      --data 'methods[]=POST' \
      --data name=auth-register

    # Shortener routes (protected)
    curl -i -X POST $KONG_ADMIN_URL/services/shortener-service/routes \
      --data 'paths[]=/api/v1/shorten' \
      --data 'methods[]=POST' \
      --data name=create-short-url

    curl -i -X POST $KONG_ADMIN_URL/services/shortener-service/routes \
      --data 'paths[]=/api/v1/urls' \
      --data 'methods[]=GET' \
      --data name=list-urls

    # Redirect routes (public)
    curl -i -X POST $KONG_ADMIN_URL/services/redirect-service/routes \
      --data 'paths[]=/r' \
      --data 'methods[]=GET' \
      --data name=redirect-short-url

    # Analytics routes (protected)
    curl -i -X POST $KONG_ADMIN_URL/services/analytics-service/routes \
      --data 'paths[]=/api/v1/analytics' \
      --data 'methods[]=GET' \
      --data name=analytics-dashboard

    # Billing routes (protected)
    curl -i -X POST $KONG_ADMIN_URL/services/billing-service/routes \
      --data 'paths[]=/api/v1/billing' \
      --data 'methods[]=GET,POST,PUT' \
      --data name=billing-routes

    echo "Routes created!"

    # Create Consumers
    echo "Creating consumers..."

    curl -i -X POST $KONG_ADMIN_URL/consumers \
      --data username=bitzap-app

    curl -i -X POST $KONG_ADMIN_URL/consumers/bitzap-app/jwt \
      --data key=bitzap-key \
      --data secret=bitzap-secret-2025

    curl -i -X POST $KONG_ADMIN_URL/consumers \
      --data username=mobile-app

    curl -i -X POST $KONG_ADMIN_URL/consumers/mobile-app/jwt \
      --data key=mobile-key \
      --data secret=mobile-secret-2025

    echo "Consumers created!"

    # Enable JWT Authentication on protected routes
    echo "� Snetting up JWT authentication..."

    # JWT for shortener routes
    curl -i -X POST $KONG_ADMIN_URL/routes/create-short-url/plugins \
      --data name=jwt \
      --data config.secret_is_base64=false

    curl -i -X POST $KONG_ADMIN_URL/routes/list-urls/plugins \
      --data name=jwt \
      --data config.secret_is_base64=false

    # JWT for analytics routes
    curl -i -X POST $KONG_ADMIN_URL/routes/analytics-dashboard/plugins \
      --data name=jwt \
      --data config.secret_is_base64=false

    # JWT for billing routes
    curl -i -X POST $KONG_ADMIN_URL/routes/billing-routes/plugins \
      --data name=jwt \
      --data config.secret_is_base64=false

    echo "JWT authentication enabled!"

    # Enable Global Plugins
    echo "Setting up global plugins..."

    # CORS
    curl -i -X POST $KONG_ADMIN_URL/plugins \
      --data name=cors \
      --data config.origins=* \
      --data config.methods=GET,POST,PUT,DELETE,OPTIONS \
      --data config.headers=Accept,Accept-Version,Content-Length,Content-MD5,Content-Type,Date,X-Auth-Token,Authorization \
      --data config.credentials=true

    # Rate Limiting
    curl -i -X POST $KONG_ADMIN_URL/plugins \
      --data name=rate-limiting \
      --data config.minute=100 \
      --data config.hour=1000 \
      --data config.day=10000

    # Request Size Limiting
    curl -i -X POST $KONG_ADMIN_URL/plugins \
      --data name=request-size-limiting \
      --data config.allowed_payload_size=10

    echo "Global plugins enabled!"
fi

echo "Kong Gateway setup completed!"
echo ""
echo "Summary:"
echo "- Kong Proxy: http://localhost:8000"
echo "- Kong Admin: http://localhost:8001"
echo "- Konga UI: http://localhost:1337"
echo ""
echo "API Endpoints:"
echo "- Auth Login: http://localhost:8000/api/v1/auth/login"
echo "- Auth Register: http://localhost:8000/api/v1/auth/register"
echo "- Create Short URL: http://localhost:8000/api/v1/shorten (JWT required)"
echo "- List URLs: http://localhost:8000/api/v1/urls (JWT required)"
echo "- Redirect: http://localhost:8000/r/{short_code}"
echo "- Analytics: http://localhost:8000/api/v1/analytics (JWT required)"
echo "- Billing: http://localhost:8000/api/v1/billing (JWT required)"
echo ""
echo "JWT Consumers:"
echo "- App: bitzap-key / bitzap-secret-2025"
echo "- Mobile: mobile-key / mobile-secret-2025"
echo ""
echo "Generate JWT token:"
echo "  Go: cd infra/kong && go run generate-jwt.go"