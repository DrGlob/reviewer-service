#!/bin/bash

echo "=== Testing Stats Endpoint ==="

# Создадим тестовые данные
echo "1. Creating test data..."
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "stats-team",
    "members": [
      {"user_id": "stat-user1", "username": "StatUser1", "is_active": true},
      {"user_id": "stat-user2", "username": "StatUser2", "is_active": true},
      {"user_id": "stat-user3", "username": "StatUser3", "is_active": true}
    ]
  }'

echo -e "\n2. Creating PRs..."
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{"pull_request_id":"stat-pr1","pull_request_name":"Stat PR1","author_id":"stat-user1"}'

curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{"pull_request_id":"stat-pr2","pull_request_name":"Stat PR2","author_id":"stat-user2"}'

echo -e "\n3. Testing reassign to create more assignments..."
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{"pull_request_id":"stat-pr1","old_user_id":"stat-user2"}'

echo -e "\n4. Getting statistics..."
curl -s http://localhost:8080/stats | python3 -m json.tool

echo -e "\n=== Stats Test Completed ==="