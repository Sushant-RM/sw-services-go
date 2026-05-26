#!/bin/bash
# verifier_flow.sh: Simulates verifier and inspector stages executing direct curl JSON calls

echo -e "\n======================================================================"
echo -e "                   STAGE 2: VERIFIER & INSPECTION FLOW"
echo -e "======================================================================"

if [ ! -f scratch/.current_app_no ]; then
  echo -e "\033[0;31m[FAIL]\033[0m No active Application Number found in scratch context. Run citizen_flow.sh first."
  exit 1
fi

APP_NO=$(cat scratch/.current_app_no)
echo -e "Processing active application: \033[1;33m$APP_NO\033[0m"

# Document Verification
echo -n "Verifying applicant documents (VERIFY_AND_FORWARD)... "
v_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
  -H "Content-Type: application/json" \
  -d "{
    \"RequestInfo\": {
      \"apiId\": \"Rainmaker\", \"ver\": \".01\", \"ts\": 1699999999000, \"authToken\": \"test-token\",
      \"userInfo\": {
        \"id\": 79, \"uuid\": \"c3a2f8c8-b046-4f68-95dc-15c9d65b8ead\", \"userName\": \"admin\",
        \"roles\": [{\"name\": \"Super User\", \"code\": \"SUPERUSER\", \"tenantId\": \"pb\"}]
      }
    },
    \"SewerageConnection\": {
      \"applicationNo\": \"$APP_NO\",
      \"tenantId\": \"pb.amritsar\",
      \"processInstance\": {
        \"action\": \"VERIFY_AND_FORWARD\"
      }
    }
  }")

status=$(echo "$v_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)

if [ "$status" != "PENDING_FOR_FIELD_INSPECTION" ]; then
  echo -e "\033[0;31m[FAIL]\033[0m Verification step failed. Status: $status"
  echo "$v_res"
  exit 1
fi
echo -e "\033[0;32m[PASS]\033[0m Advanced state to: \033[1;32m$status\033[0m"

# Field Inspection
echo -n "Performing physical site feasibility inspection (VERIFY_AND_FORWARD)... "
i_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
  -H "Content-Type: application/json" \
  -d "{
    \"RequestInfo\": {
      \"apiId\": \"Rainmaker\", \"ver\": \".01\", \"ts\": 1699999999000, \"authToken\": \"test-token\",
      \"userInfo\": {
        \"id\": 79, \"uuid\": \"c3a2f8c8-b046-4f68-95dc-15c9d65b8ead\", \"userName\": \"admin\",
        \"roles\": [{\"name\": \"Super User\", \"code\": \"SUPERUSER\", \"tenantId\": \"pb\"}]
      }
    },
    \"SewerageConnection\": {
      \"applicationNo\": \"$APP_NO\",
      \"tenantId\": \"pb.amritsar\",
      \"processInstance\": {
        \"action\": \"VERIFY_AND_FORWARD\"
      }
    }
  }")

status=$(echo "$i_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)

if [ "$status" != "PENDING_APPROVAL_FOR_CONNECTION" ]; then
  echo -e "\033[0;31m[FAIL]\033[0m Physical inspection step failed. Status: $status"
  echo "$i_res"
  exit 1
fi
echo -e "\033[0;32m[PASS]\033[0m Advanced state to: \033[1;32m$status\033[0m"
exit 0
