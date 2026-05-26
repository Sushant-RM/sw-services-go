#!/bin/bash
# citizen_flow.sh: Simulates citizen application registration and initial submission

echo -e "\n======================================================================"
echo -e "                   STAGE 1: CITIZEN REGISTRATION FLOW"
echo -e "======================================================================"

# Ingest application
echo -n "Registering new connection application... "
create_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_create" \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {
      "apiId": "Rainmaker",
      "ver": ".01",
      "ts": 1699999999000,
      "authToken": "test-token",
      "userInfo": {
        "id": 79,
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

APP_NO=$(echo "$create_res" | jq -r '.SewerageConnections[0].applicationNo' 2>/dev/null)

if [ -z "$APP_NO" ] || [ "$APP_NO" == "null" ]; then
  echo -e "\033[0;31m[FAIL]\033[0m Ingestion failed or response payload malformed."
  echo "$create_res"
  exit 1
fi

echo -e "\033[0;32m[PASS]\033[0m Application created: \033[1;33m$APP_NO\033[0m"
echo "$APP_NO" > scratch/.current_app_no

# Submit application
echo -n "Submitting application for document verification... "
submit_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
  -H "Content-Type: application/json" \
  -d "{
    \"RequestInfo\": {
      \"apiId\": \"Rainmaker\", \"ver\": \".01\", \"ts\": 1699999999000, \"authToken\": \"test-token\",
      \"userInfo\": {
        \"id\": 79, \"uuid\": \"c3a2f8c8-b046-4f68-95dc-15c9d65b8ead\", \"userName\": \"admin\",
        \"roles\": [{\"name\": \"Super User\", \"code\": \"SUPERUSER\", \"tenantId\": \"pb\"}], \"tenantId\": \"pb\"
      }
    },
    \"SewerageConnection\": {
      \"applicationNo\": \"$APP_NO\",
      \"tenantId\": \"pb.amritsar\",
      \"propertyId\": \"PT-107-123456\",
      \"connectionType\": \"Non Metered\",
      \"processInstance\": {
        \"action\": \"SUBMIT_APPLICATION\"
      }
    }
  }")

status=$(echo "$submit_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)

if [ "$status" == "PENDING_FOR_DOCUMENT_VERIFICATION" ]; then
  echo -e "\033[0;32m[PASS]\033[0m Advanced state to: \033[1;32m$status\033[0m"
  exit 0
else
  echo -e "\033[0;31m[FAIL]\033[0m Workflow rejection. Status: $status"
  echo "$submit_res"
  exit 1
fi
