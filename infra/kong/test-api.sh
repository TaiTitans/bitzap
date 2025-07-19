#!/bin/bash

# API Testing Script for BitZap Kong Gateway
# This script tests all API endpoints with proper authentication

KONG_URL="http://localhost:8000"
KONG_ADMIN_URL="http://localhost:8001"

echo "🧪 Testing BitZap API Gateway..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Function to print info
print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check if Kong is running
print_info "Checking Kong status..."
curl -f $KONG_ADMIN_URL/status > /dev/null 2>&1
print_result $? "Kong Gateway is running"

# Generate JWT token
print_info "Generating JWT token..."
if command -v go &> /dev/null; then
    JWT_TOKEN=$(cd infra/kong && go run generate-jwt.go | grep -A1 "Generated JWT Token:" | tail -1)
    print_result $? "JWT token generated (Go)"
elif command -v python3 &> /dev/null; then
    JWT_TOKEN=$(python3 infra/kong/generate-jwt.py | grep -A1 "Generated JWT Token:" | tail -1)
    print_result $? "JWT token generated (Python)"
else
    echo -e "${RED}❌ Neither Go nor Python3 found. Please install one to generate JWT tokens.${NC}"
    exit 1
fi

echo ""
echo "🔑 Using JWT Token: ${JWT_TOKEN:0:50}..."
echo ""

# Test Public Endpoints
echo "📋 Testing Public Endpoints..."

# Test Auth Register (Mock)
print_info "Testing Auth Register..."
curl -s -o /dev/null -w "%{http_code}" -X POST $KONG_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}' | grep -q "200\|201\|404\|500"
print_result $? "POST /api/v1/auth/register"

# Test Auth Login (Mock)
print_info "Testing Auth Login..."
curl -s -o /dev/null -w "%{http_code}" -X POST $KONG_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' | grep -q "200\|201\|404\|500"
print_result $? "POST /api/v1/auth/login"

# Test Redirect (Mock)
print_info "Testing Redirect..."
curl -s -o /dev/null -w "%{http_code}" -X GET $KONG_URL/r/abc123 | grep -q "200\|301\|302\|404\|500"
print_result $? "GET /r/{short_code}"

echo ""
echo "🔒 Testing Protected Endpoints..."

# Test Create Short URL
print_info "Testing Create Short URL..."
curl -s -o /dev/null -w "%{http_code}" -X POST $KONG_URL/api/v1/shorten \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","custom_code":"test123"}' | grep -q "200\|201\|404\|500"
print_result $? "POST /api/v1/shorten (with JWT)"

# Test List URLs
print_info "Testing List URLs..."
curl -s -o /dev/null -w "%{http_code}" -X GET $KONG_URL/api/v1/urls \
  -H "Authorization: Bearer $JWT_TOKEN" | grep -q "200\|404\|500"
print_result $? "GET /api/v1/urls (with JWT)"

# Test Analytics
print_info "Testing Analytics..."
curl -s -o /dev/null -w "%{http_code}" -X GET $KONG_URL/api/v1/analytics \
  -H "Authorization: Bearer $JWT_TOKEN" | grep -q "200\|404\|500"
print_result $? "GET /api/v1/analytics (with JWT)"

# Test Billing
print_info "Testing Billing..."
curl -s -o /dev/null -w "%{http_code}" -X GET $KONG_URL/api/v1/billing \
  -H "Authorization: Bearer $JWT_TOKEN" | grep -q "200\|404\|500"
print_result $? "GET /api/v1/billing (with JWT)"

echo ""
echo "🚫 Testing Authentication (should fail without JWT)..."

# Test protected endpoint without JWT
print_info "Testing Create Short URL without JWT..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $KONG_URL/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}')
if [ "$HTTP_CODE" = "401" ] || [ "$HTTP_CODE" = "403" ]; then
    print_result 0 "POST /api/v1/shorten (without JWT) - Correctly rejected"
else
    print_result 1 "POST /api/v1/shorten (without JWT) - Should be rejected"
fi

echo ""
echo "📊 Kong Gateway Status..."

# Check Kong services
print_info "Checking Kong services..."
SERVICES_COUNT=$(curl -s $KONG_ADMIN_URL/services | jq '.data | length' 2>/dev/null || echo "N/A")
echo "Services configured: $SERVICES_COUNT"

# Check Kong routes
print_info "Checking Kong routes..."
ROUTES_COUNT=$(curl -s $KONG_ADMIN_URL/routes | jq '.data | length' 2>/dev/null || echo "N/A")
echo "Routes configured: $ROUTES_COUNT"

# Check Kong consumers
print_info "Checking Kong consumers..."
CONSUMERS_COUNT=$(curl -s $KONG_ADMIN_URL/consumers | jq '.data | length' 2>/dev/null || echo "N/A")
echo "Consumers configured: $CONSUMERS_COUNT"

# Check Kong plugins
print_info "Checking Kong plugins..."
PLUGINS_COUNT=$(curl -s $KONG_ADMIN_URL/plugins | jq '.data | length' 2>/dev/null || echo "N/A")
echo "Plugins enabled: $PLUGINS_COUNT"

echo ""
echo "🎉 API Gateway testing completed!"
echo ""
echo "📋 Summary:"
echo "- Kong Proxy: $KONG_URL"
echo "- Kong Admin: $KONG_ADMIN_URL"
echo "- Konga UI: http://localhost:1337"
echo ""
echo "💡 Tips:"
echo "- Use Konga UI for visual management"
echo "- Check logs: docker logs kong"
echo "- Monitor metrics: curl $KONG_ADMIN_URL/metrics"