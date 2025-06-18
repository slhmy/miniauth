#!/bin/bash

# Internal OAuth API Test Script with Internal Token Support
# This script demonstrates how to use the internal OAuth API with internal token authentication

BASE_URL="http://localhost:8080/api"
INTERNAL_TOKEN="${MINIAUTH_INTERNAL_TOKEN:-miniauth-internal-default-token-change-in-production}"
ADMIN_SESSION=""  # Optional: admin session cookie as fallback

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Internal OAuth API Test Script with Internal Token ===${NC}"
echo ""
echo -e "${BLUE}Using internal token: ${INTERNAL_TOKEN}${NC}"
echo ""

echo -e "${YELLOW}Testing single application creation with different authentication methods...${NC}"

# Test 1: Create with Authorization header (Bearer)
echo -e "${GREEN}1. Creating application with Authorization: Bearer header...${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $INTERNAL_TOKEN" \
    -d '{
        "name": "Test App (Bearer Auth)",
        "description": "Test application created via Bearer token",
        "website": "https://test-bearer.example.com",
        "redirect_uris": ["https://test-bearer.example.com/oauth/callback"],
        "scopes": ["read", "profile"],
        "trusted": true,
        "client_id": "test-bearer-app-123",
        "client_secret": "test-bearer-secret-456"
    }')

echo "Response: $RESPONSE"
echo ""

# Test 2: Create with Authorization header (Internal)
echo -e "${GREEN}2. Creating application with Authorization: Internal header...${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Internal $INTERNAL_TOKEN" \
    -d '{
        "name": "Test App (Internal Auth)",
        "description": "Test application created via Internal token",
        "website": "https://test-internal.example.com",
        "redirect_uris": ["https://test-internal.example.com/oauth/callback"],
        "scopes": ["read", "profile"],
        "trusted": true,
        "client_id": "test-internal-app-123",
        "client_secret": "test-internal-secret-456"
    }')

echo "Response: $RESPONSE"
echo ""

# Test 3: Create with X-Internal-Token header
echo -e "${GREEN}3. Creating application with X-Internal-Token header...${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications" \
    -H "Content-Type: application/json" \
    -H "X-Internal-Token: $INTERNAL_TOKEN" \
    -d '{
        "name": "Test App (X-Internal-Token)",
        "description": "Test application created via X-Internal-Token header",
        "website": "https://test-header.example.com",
        "redirect_uris": ["https://test-header.example.com/oauth/callback"],
        "scopes": ["read"],
        "trusted": false,
        "client_id": "test-header-app-123",
        "client_secret": "test-header-secret-456"
    }')

echo "Response: $RESPONSE"
echo ""

# Test 4: Create with query parameter
echo -e "${GREEN}4. Creating application with internal_token query parameter...${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications?internal_token=$INTERNAL_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Test App (Query Param)",
        "description": "Test application created via query parameter",
        "website": "https://test-query.example.com",
        "redirect_uris": ["https://test-query.example.com/oauth/callback"],
        "scopes": ["read"],
        "trusted": false,
        "client_id": "test-query-app-123",
        "client_secret": "test-query-secret-456"
    }')

echo "Response: $RESPONSE"
echo ""

# Test 5: Try to create with duplicate client_id (should fail)
echo -e "${YELLOW}5. Testing duplicate client_id (should fail)...${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $INTERNAL_TOKEN" \
    -d '{
        "name": "Duplicate App",
        "description": "This should fail due to duplicate client_id",
        "redirect_uris": ["https://duplicate.example.com/callback"],
        "client_id": "test-bearer-app-123",
        "client_secret": "different-secret"
    }')

echo "Response: $RESPONSE"
echo ""

# Test 6: Batch creation with internal token
echo -e "${YELLOW}6. Testing batch application creation with internal token...${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications/batch" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $INTERNAL_TOKEN" \
    -d '[
        {
            "name": "Batch App 1 (Token)",
            "description": "First batch application via token",
            "redirect_uris": ["https://batch1-token.example.com/callback"],
            "client_id": "batch-token-app-1",
            "client_secret": "batch-token-secret-1"
        },
        {
            "name": "Batch App 2 (Token)", 
            "description": "Second batch application via token",
            "redirect_uris": ["https://batch2-token.example.com/callback"],
            "scopes": ["read"],
            "trusted": true,
            "client_id": "batch-token-app-2",
            "client_secret": "batch-token-secret-2"
        },
        {
            "name": "Batch App 3 (duplicate)",
            "description": "This should fail due to duplicate client_id",
            "redirect_uris": ["https://batch3-token.example.com/callback"],
            "client_id": "test-bearer-app-123",
            "client_secret": "batch-token-secret-3"
        }
    ]')

echo "Response: $RESPONSE"
echo ""

# Test 7: Test unauthorized access (invalid token)
echo -e "${YELLOW}7. Testing unauthorized access with invalid token...${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer invalid-token" \
    -d '{
        "name": "Unauthorized App",
        "description": "This should fail",
        "redirect_uris": ["https://unauthorized.example.com/callback"],
        "client_id": "unauthorized-app",
        "client_secret": "unauthorized-secret"
    }')

echo "Response: $RESPONSE"
echo ""

# Test 8: Test session fallback (if admin session is provided)
if [ -n "$ADMIN_SESSION" ]; then
    echo -e "${YELLOW}8. Testing admin session fallback...${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/admin/oauth/internal/applications" \
        -H "Content-Type: application/json" \
        -H "Cookie: $ADMIN_SESSION" \
        -d '{
            "name": "Session Fallback App",
            "description": "Created via admin session fallback",
            "redirect_uris": ["https://session.example.com/callback"],
            "client_id": "session-app-123",
            "client_secret": "session-secret-456"
        }')
    
    echo "Response: $RESPONSE"
    echo ""
fi

echo -e "${GREEN}Test completed!${NC}"
echo ""
echo -e "${YELLOW}Authentication methods tested:${NC}"
echo "✓ Authorization: Bearer <token>"
echo "✓ Authorization: Internal <token>"
echo "✓ X-Internal-Token: <token>"
echo "✓ ?internal_token=<token>"
if [ -n "$ADMIN_SESSION" ]; then
    echo "✓ Admin session fallback"
fi
echo ""
echo -e "${YELLOW}To run this script:${NC}"
echo "1. Start the miniauth server: make run"
echo "2. Set internal token (optional): export MINIAUTH_INTERNAL_TOKEN='your-secret-token'"
echo "3. Run this script: bash examples/test_internal_oauth_api.sh"
echo ""
echo -e "${BLUE}Default internal token: miniauth-internal-default-token-change-in-production${NC}"
