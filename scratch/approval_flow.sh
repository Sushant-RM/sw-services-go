#!/bin/bash
# approval_flow.sh: Simulates officer approvals, payments, and live activation executing direct curls

echo -e "\n======================================================================"
echo -e "                   STAGE 3: APPROVAL, BILLING & ACTIVATION"
echo -e "======================================================================"

if [ ! -f scratch/.current_app_no ]; then
  echo -e "\033[0;31m[FAIL]\033[0m No active Application Number found in scratch context."
  exit 1
fi

APP_NO=$(cat scratch/.current_app_no)
echo -e "Processing active application: \033[1;33m$APP_NO\033[0m"

# Approval Officer
echo -n "Approving application and generating Connection Number... "
a_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
        \"action\": \"APPROVE_FOR_CONNECTION\"
      }
    }
  }")

status=$(echo "$a_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
conn_no=$(echo "$a_res" | jq -r '.SewerageConnections[0].connectionNo' 2>/dev/null)

if [ "$status" != "PENDING_FOR_PAYMENT" ]; then
  echo -e "\033[0;31m[FAIL]\033[0m Connection approval failed. Status: $status"
  echo "$a_res"
  exit 1
fi
echo -e "\033[0;32m[PASS]\033[0m Approved! Legal Connection ID allocated: \033[1;36m$conn_no\033[0m"

# Billing Settle Payment
echo -n "Apportioning tax ledgers & processing payment fee... "
p_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
        \"action\": \"PAY\"
      }
    }
  }")

status=$(echo "$p_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)

if [ "$status" != "PENDING_FOR_CONNECTION_ACTIVATION" ]; then
  echo -e "\033[0;31m[FAIL]\033[0m Payment processing failed. Status: $status"
  echo "$p_res"
  exit 1
fi
echo -e "\033[0;32m[PASS]\033[0m Payment transaction settled successfully."

# Connection Activation
echo -n "Activating sewerage connection line to live operational state... "
act_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
        \"action\": \"ACTIVATE_CONNECTION\"
      }
    }
  }")

status=$(echo "$act_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)

if [ "$status" != "CONNECTION_ACTIVATED" ]; then
  echo -e "\033[0;31m[FAIL]\033[0m Activation failed. Status: $status"
  echo "$act_res"
  exit 1
fi
echo -e "\033[0;32m[PASS]\033[0m Sewerage Connection is now officially \033[1;32mACTIVATED\033[0m."
exit 0
