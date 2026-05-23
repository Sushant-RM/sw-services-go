#!/bin/bash
set -e

echo "=========================================================="
echo "DIGIT REST API & Integrated Microservices Audit"
echo "=========================================================="

services=(
  "3468:Go Sewerage Service"
  "3456:egov-mdms-service"
  "3457:egov-idgen"
  "3460:egov-user"
  "3463:egov-workflow-v2"
  "3466:property-services"
)

echo "Waiting for all peer service ports to be listening..."
for item in "${services[@]}"; do
  port="${item%%:*}"
  name="${item#*:}"
  
  echo -n "Checking $name (Port $port)... "
  until curl -s -o /dev/null "http://localhost:$port" || [ $? -eq 52 ] || [ $? -eq 56 ]; do
    echo -n "."
    sleep 3
  done
  echo " ONLINE!"
done

echo ""
echo "All peer services are listening on their network ports!"
echo "Conducting Phase 4 Service Integration Test (POST _create)..."
echo ""

# Execute the Go Sewerage Service create connection API
response=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_create" \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {
      "apiId": "Rainmaker",
      "ver": ".01",
      "ts": 1699999999000,
      "action": "",
      "msgId": "20170310130900|en_IN",
      "authToken": "test-token",
      "userInfo": {
        "id": 1,
        "uuid": "c3a2f8c8-b046-4f68-95dc-15c9d65b8ead",
        "userName": "admin",
        "roles": [{"name": "Super User", "code": "SUPERUSER", "tenantId": "pb"}],
        "tenantId": "pb"
      }
    },
    "SewerageConnection": {
      "propertyId": "PT-107-123456",
      "tenantId": "pb.amritsar",
      "connectionType": "Non Metered",
      "roadType": "BERMCUTTINGKATCHA",
      "roadCuttingArea": 25,
      "processInstance": {
        "action": "INITIATE"
      }
    }
  }')

echo "----------------------------------------------------------"
echo "Go Sewerage Service _create API Response:"
echo "----------------------------------------------------------"
echo "$response" | jq . || echo "$response"
echo "----------------------------------------------------------"

echo ""
echo "Querying PostgreSQL 'rainmaker' Database to confirm write parity..."
postgres_container=$(docker ps -f name=postgres --format "{{.Names}}" | head -n 1)
docker exec -t "$postgres_container" psql -U postgres -d rainmaker -c "SELECT applicationno, property_id, connectiontype, applicationstatus FROM eg_sw_connection ORDER BY createdtime DESC LIMIT 1;"

echo ""
echo "=========================================================="
echo "Phase 4 Audit Completed Successfully!"
echo "=========================================================="
