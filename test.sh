#!/bin/bash

SUFFIX=$(date +%s)

echo "=== PR Reassign Test (Run $SUFFIX) ==="
echo

# 1. Создаем команду
echo "1. Creating team..."
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "team-'$SUFFIX'",
    "members": [
      {"user_id": "author-'$SUFFIX'", "username": "Author", "is_active": true},
      {"user_id": "reviewer1-'$SUFFIX'", "username": "Reviewer1", "is_active": true},
      {"user_id": "reviewer2-'$SUFFIX'", "username": "Reviewer2", "is_active": true},
      {"user_id": "reviewer3-'$SUFFIX'", "username": "Reviewer3", "is_active": true}
    ]
  }'
echo
echo

# 2. Создаем PR
echo "2. Creating PR..."
RESPONSE=$(curl -s -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-'$SUFFIX'",
    "pull_request_name": "Test Feature",
    "author_id": "author-'$SUFFIX'"
  }')
echo "$RESPONSE"
echo

# 3. Извлекаем назначенных ревьюеров из ответа
REVIEWER1=$(echo "$RESPONSE" | grep -o '"assigned_reviewers":\["[^"]*"' | cut -d'"' -f4 | head -1)
REVIEWER2=$(echo "$RESPONSE" | grep -o '"assigned_reviewers":\["[^"]*","[^"]*"' | cut -d'"' -f6 | head -1)

echo "3. Detected reviewers: $REVIEWER1 and $REVIEWER2"
echo

# 4. Проверяем назначенных ревьюеров
echo "4. Checking initial reviewers..."
echo "$REVIEWER1 review PRs:"
curl -s "http://localhost:8080/users/getReview?user_id=$REVIEWER1" | python3 -m json.tool
echo

echo "$REVIEWER2 review PRs:"
curl -s "http://localhost:8080/users/getReview?user_id=$REVIEWER2" | python3 -m json.tool
echo

# 5. Тестируем reassign (заменяем первого ревьюера)
echo "5. Executing reassign (replacing $REVIEWER1)..."
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-'$SUFFIX'",
    "old_user_id": "'$REVIEWER1'"
  }'
echo
echo

# 6. Проверяем результат
echo "6. Checking reassign result..."
echo "$REVIEWER1 review PRs (should be empty):"
curl -s "http://localhost:8080/users/getReview?user_id=$REVIEWER1" | python3 -m json.tool
echo

echo "=== Test completed ==="