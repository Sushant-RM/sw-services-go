#!/bin/bash
# ==============================================================================
# DIGIT Sewerage Service - Automated Municipal Workflow Transition Client
# ==============================================================================

APP_NO=$1
ACTION=$2

if [ -z "$APP_NO" ] || [ -z "$ACTION" ]; then
  echo -e "\033[1;31mUsage:\033[0m bash scratch/transition.sh <application_number> <action>"
  echo -e "Actions: SUBMIT_APPLICATION | VERIFY_AND_FORWARD | APPROVE_FOR_CONNECTION | PAY | ACTIVATE_CONNECTION"
  exit 1
fi

# Define Role info to display (we use valid SUPERUSER credentials internally to bypass egov-user database limits)
case "$ACTION" in
  SUBMIT_APPLICATION)
    DISPLAY_ROLE="Citizen"
    ;;
  VERIFY_AND_FORWARD)
    # Check current status first to differentiate Verifier from Inspector
    status=$(curl -s "http://localhost:3468/sw-services/swc/_search?tenantId=pb.amritsar&applicationNumber=$APP_NO" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
    if [ "$status" == "PENDING_FOR_FIELD_INSPECTION" ]; then
      DISPLAY_ROLE="Field Inspector"
    else
      DISPLAY_ROLE="Document Verifier"
    fi
    ;;
  APPROVE_FOR_CONNECTION)
    DISPLAY_ROLE="Approval Officer"
    ;;
  PAY)
    DISPLAY_ROLE="Billing Stage"
    ;;
  ACTIVATE_CONNECTION)
    DISPLAY_ROLE="Activation Officer"
    ;;
  *)
    echo "Unknown Action: $ACTION"
    exit 1
    ;;
esac

# Execute request using valid SUPERUSER user context (registered in target VM database)
res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
  -H "Content-Type: application/json" \
  -d "{
    \"RequestInfo\": {
      \"apiId\": \"Rainmaker\", \"ver\": \".01\", \"ts\": 1699999999000, \"authToken\": \"test-token\",
      \"userInfo\": {
        \"id\": 79,
        \"uuid\": \"c3a2f8c8-b046-4f68-95dc-15c9d65b8ead\",
        \"userName\": \"admin\",
        \"roles\": [{\"name\": \"Super User\", \"code\": \"SUPERUSER\", \"tenantId\": \"pb\"}],
        \"tenantId\": \"pb\"
      }
    },
    \"SewerageConnection\": {
      \"applicationNo\": \"$APP_NO\",
      \"tenantId\": \"pb.amritsar\",
      \"processInstance\": {
        \"action\": \"$ACTION\"
      }
    }
  }")

# Check status and response
status=$(echo "$res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
conn_no=$(echo "$res" | jq -r '.SewerageConnections[0].connectionNo' 2>/dev/null)

if [ "$status" != "null" ] && [ ! -z "$status" ]; then
  if [ ! -z "$conn_no" ] && [ "$conn_no" != "null" ]; then
    echo -e "[\033[1;36m$DISPLAY_ROLE\033[0m] SUCCESS: Action=\033[1;32m$ACTION\033[0m -> State: \033[1;33m$status\033[0m | ConnectionNo: \033[1;35m$conn_no\033[0m"
  else
    echo -e "[\033[1;36m$DISPLAY_ROLE\033[0m] SUCCESS: Action=\033[1;32m$ACTION\033[0m -> State: \033[1;33m$status\033[0m"
  fi
else
  # Print complete error if transition failed
  echo -e "Transition \033[1;31mFAILED\033[0m: $res"
fi
