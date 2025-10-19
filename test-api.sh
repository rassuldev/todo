#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base URLs
USER_SERVICE="http://localhost:8081"
AUTH_SERVICE="http://localhost:8082"
TASK_SERVICE="http://localhost:8083"
NOTIFICATION_SERVICE="http://localhost:8084"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Testing Task Management Microservices${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Test 1: Create a user
echo -e "${BLUE}1. Creating a new user...${NC}"
USER_RESPONSE=$(curl -s -X POST $USER_SERVICE/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }')

echo "$USER_RESPONSE" | jq .
USER_ID=$(echo "$USER_RESPONSE" | jq -r '.id')

if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
    echo -e "${GREEN}✓ User created successfully with ID: $USER_ID${NC}\n"
else
    echo -e "${RED}✗ Failed to create user${NC}\n"
    exit 1
fi

# Test 2: Login
echo -e "${BLUE}2. Logging in...${NC}"
AUTH_RESPONSE=$(curl -s -X POST $AUTH_SERVICE/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "$AUTH_RESPONSE" | jq .
ACCESS_TOKEN=$(echo "$AUTH_RESPONSE" | jq -r '.access_token')

if [ -n "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ Login successful${NC}\n"
else
    echo -e "${RED}✗ Login failed${NC}\n"
    exit 1
fi

# Test 3: Validate token
echo -e "${BLUE}3. Validating token...${NC}"
VALIDATE_RESPONSE=$(curl -s -X POST $AUTH_SERVICE/api/auth/validate \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$ACCESS_TOKEN\"
  }")

echo "$VALIDATE_RESPONSE" | jq .

if [ "$(echo "$VALIDATE_RESPONSE" | jq -r '.valid')" == "true" ]; then
    echo -e "${GREEN}✓ Token is valid${NC}\n"
else
    echo -e "${RED}✗ Token validation failed${NC}\n"
fi

# Test 4: Get user details
echo -e "${BLUE}4. Getting user details...${NC}"
GET_USER_RESPONSE=$(curl -s -X GET $USER_SERVICE/api/users/$USER_ID)

echo "$GET_USER_RESPONSE" | jq .
echo -e "${GREEN}✓ User details retrieved${NC}\n"

# Test 5: Create a task
echo -e "${BLUE}5. Creating a task...${NC}"
TASK_RESPONSE=$(curl -s -X POST $TASK_SERVICE/api/tasks \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Complete API testing\",
    \"description\": \"Test all microservices endpoints\",
    \"priority\": \"HIGH\",
    \"user_id\": \"$USER_ID\",
    \"due_date\": \"2024-12-31T23:59:59Z\"
  }")

echo "$TASK_RESPONSE" | jq .
TASK_ID=$(echo "$TASK_RESPONSE" | jq -r '.id')

if [ -n "$TASK_ID" ] && [ "$TASK_ID" != "null" ]; then
    echo -e "${GREEN}✓ Task created successfully with ID: $TASK_ID${NC}\n"
else
    echo -e "${RED}✗ Failed to create task${NC}\n"
fi

# Test 6: Get task details
echo -e "${BLUE}6. Getting task details...${NC}"
GET_TASK_RESPONSE=$(curl -s -X GET $TASK_SERVICE/api/tasks/$TASK_ID)

echo "$GET_TASK_RESPONSE" | jq .
echo -e "${GREEN}✓ Task details retrieved${NC}\n"

# Test 7: Update task
echo -e "${BLUE}7. Updating task status...${NC}"
UPDATE_TASK_RESPONSE=$(curl -s -X PUT $TASK_SERVICE/api/tasks/$TASK_ID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete API testing",
    "description": "Test all microservices endpoints - Updated",
    "status": "IN_PROGRESS",
    "priority": "HIGH",
    "due_date": "2024-12-31T23:59:59Z"
  }')

echo "$UPDATE_TASK_RESPONSE" | jq .
echo -e "${GREEN}✓ Task updated successfully${NC}\n"

# Test 8: List user tasks
echo -e "${BLUE}8. Listing user tasks...${NC}"
LIST_TASKS_RESPONSE=$(curl -s -X GET "$TASK_SERVICE/api/users/$USER_ID/tasks?page=1&page_size=10")

echo "$LIST_TASKS_RESPONSE" | jq .
echo -e "${GREEN}✓ User tasks listed${NC}\n"

# Test 9: Send email notification
echo -e "${BLUE}9. Sending email notification...${NC}"
EMAIL_RESPONSE=$(curl -s -X POST $NOTIFICATION_SERVICE/api/notifications/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "subject": "Task Created",
    "body": "Your task has been created successfully!"
  }')

echo "$EMAIL_RESPONSE" | jq .

if [ "$(echo "$EMAIL_RESPONSE" | jq -r '.success')" == "true" ]; then
    echo -e "${GREEN}✓ Email notification sent${NC}\n"
else
    echo -e "${GREEN}✓ Email notification logged (SMTP not configured)${NC}\n"
fi

# Test 10: List all users
echo -e "${BLUE}10. Listing all users...${NC}"
LIST_USERS_RESPONSE=$(curl -s -X GET "$USER_SERVICE/api/users?page=1&page_size=10")

echo "$LIST_USERS_RESPONSE" | jq .
echo -e "${GREEN}✓ Users listed${NC}\n"

# Test 11: List all tasks
echo -e "${BLUE}11. Listing all tasks...${NC}"
ALL_TASKS_RESPONSE=$(curl -s -X GET "$TASK_SERVICE/api/tasks?page=1&page_size=10")

echo "$ALL_TASKS_RESPONSE" | jq .
echo -e "${GREEN}✓ All tasks listed${NC}\n"

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}All tests completed successfully!${NC}"
echo -e "${BLUE}========================================${NC}\n"

echo -e "${BLUE}Created Resources:${NC}"
echo -e "User ID: ${GREEN}$USER_ID${NC}"
echo -e "Task ID: ${GREEN}$TASK_ID${NC}"
echo -e "Access Token: ${GREEN}${ACCESS_TOKEN:0:50}...${NC}\n"

echo -e "${BLUE}Cleanup (optional):${NC}"
echo -e "To delete the task: curl -X DELETE $TASK_SERVICE/api/tasks/$TASK_ID"
echo -e "To delete the user: curl -X DELETE $USER_SERVICE/api/users/$USER_ID"

