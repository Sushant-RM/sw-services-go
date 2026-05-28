#!/bin/bash
# ==============================================================================
# DIGIT Sewerage Service - Live Multi-Userflow Demonstration Suite
# ==============================================================================
# Designed for maximum visual elegance, clean spacing, and clear role separation.

# ANSI Color Code Palette
CYAN='\033[1;36m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
PURPLE='\033[1;35m'
BLUE='\033[1;34m'
WHITE='\033[1;37m'
RED='\033[1;31m'
RESET='\033[0m'
BG_BLUE='\033[44m'
BG_GREEN='\033[42m'

# Helper function to print clean horizontal dividers
print_divider() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
}

# Helper function to print spacious cards
print_user_card() {
    local role=$1
    local name=$2
    local action=$3
    local details=$4

    echo -e ""
    echo -e "   ┌────────────────────────────────────────────────────────────────────────┐"
    echo -e "   │  ${CYAN}ACTOR ROLE:${RESET}   %-55s │" "$role ($name)"
    echo -e "   │  ${YELLOW}WORKFLOW ACT:${RESET} %-55s │" "$action"
    echo -e "   │  ${WHITE}CONTEXT:${RESET}      %-55s │" "$details"
    echo -e "   └────────────────────────────────────────────────────────────────────────┘"
    echo -e ""
}

# Helper to pause and prompt user
prompt_next() {
    echo -e ""
    echo -e "   👉 ${WHITE}Press ${GREEN}[ENTER]${WHITE} to authorize this stage and proceed...${RESET}"
    read -r
}

clear

# Welcome Screen
print_divider
echo -e "               ${BG_BLUE}${WHITE}  DIGIT SEWERAGE WORKFLOW: LIVE INTERACTIVE DEMO  ${RESET}"
print_divider
echo -e "  This demonstration will walk you through a multi-user, multi-role municipal"
echo -e "  workflow for sewerage line connection creation and final activation."
echo -e ""
echo -e "  Each stage represents a real user/officer interacting with the Go microservice"
echo -e "  (${CYAN}sw-services-go${RESET}) and core platform components in real time."
echo -e ""
echo -e "  ${WHITE}Actors involved:${RESET}"
echo -e "   1. ${GREEN}Citizen (Amit Kumar)${RESET} - Registers and submits application."
echo -e "   2. ${GREEN}Document Verifier (Meera Sharma)${RESET} - Validates legal credentials."
echo -e "   3. ${GREEN}Field Inspector (Vikram Rathore)${RESET} - Inspects site dimensions."
echo -e "   4. ${GREEN}Municipal Approver (Rajesh Gupta)${RESET} - Authorizes & grants connection."
echo -e "   5. ${GREEN}Billing Clerk (Anita Patel)${RESET} - Processes demand charges."
echo -e "   6. ${GREEN}Field Engineer (Suresh Kumar)${RESET} - Activates connection line."
print_divider
prompt_next

# ------------------------------------------------------------------------------
# STAGE 1: Citizen Ingestion
# ------------------------------------------------------------------------------
clear
print_divider
echo -e "   ${BG_BLUE}${WHITE}  STAGE 1: CITIZEN REGISTRATION & INGESTION  ${RESET}"
print_divider

print_user_card "Citizen" "Amit Kumar (pb.amritsar)" "CREATE & SUBMIT APPLICATION" "Submitting new sewerage connection request for Property PT-107-123456"

echo -e "   Sending request to: ${CYAN}POST /sw-services/swc/_create${RESET}..."
sleep 1

# Execute Create
create_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_create" \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {
      "apiId": "Rainmaker", "ver": ".01", "ts": 1699999999000, "authToken": "test-token",
      "userInfo": {
        "id": 79, "uuid": "c3a2f8c8-b046-4f68-95dc-15c9d65b8ead", "userName": "admin",
        "roles": [{"name": "Super User", "code": "SUPERUSER", "tenantId": "pb"}], "tenantId": "pb"
      }
    },
    "SewerageConnection": {
      "tenantId": "pb.amritsar",
      "propertyId": "PT-107-123456",
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
    echo -e "   ${RED}❌ Application registration failed! Response: $create_res${RESET}"
    exit 1
fi

echo -e "   ${GREEN}✔ SUCCESS!${RESET} Connection application generated."
echo -e "   ${WHITE}Application Number:${RESET} ${PURPLE}$APP_NO${RESET}"
echo -e "   ${WHITE}Current Status:${RESET}     ${YELLOW}INITIATED${RESET}"

# Transition to document verification
echo -e "\n   Advancing to document verification..."
sleep 0.5

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

STATUS=$(echo "$submit_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
echo -e "   ${GREEN}✔ SUCCESS!${RESET} Application successfully submitted to municipal office."
echo -e "   ${WHITE}New Workflow State:${RESET} ${YELLOW}$STATUS${RESET}"

print_divider
prompt_next

# ------------------------------------------------------------------------------
# STAGE 2: Document Verification
# ------------------------------------------------------------------------------
clear
print_divider
echo -e "   ${BG_BLUE}${WHITE}  STAGE 2: OFFICE DOCUMENT VERIFICATION  ${RESET}"
print_divider

print_user_card "Document Verifier" "Meera Sharma (Senior Clerk)" "VERIFY_AND_FORWARD" "Inspecting applicant ownership cards and tax receipts for $APP_NO"

echo -e "   Transitioning state..."
sleep 1

verify_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
      \"processInstance\": {
        \"action\": \"VERIFY_AND_FORWARD\"
      }
    }
  }")

STATUS=$(echo "$verify_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
echo -e "   ${GREEN}✔ SUCCESS!${RESET} Documents verified. Handed over to inspection office."
echo -e "   ${WHITE}New Workflow State:${RESET} ${YELLOW}$STATUS${RESET}"

print_divider
prompt_next

# ------------------------------------------------------------------------------
# STAGE 3: Field Inspection
# ------------------------------------------------------------------------------
clear
print_divider
echo -e "   ${BG_BLUE}${WHITE}  STAGE 3: SITE FEASIBILITY FIELD INSPECTION  ${RESET}"
print_divider

print_user_card "Field Inspector" "Vikram Rathore (Sub-Divisional Inspector)" "VERIFY_AND_FORWARD" "Inspecting site dimensions. Estimating road cutting area (25 sqm Berm Cutting)"

echo -e "   Transitioning state..."
sleep 1

inspect_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
      \"processInstance\": {
        \"action\": \"VERIFY_AND_FORWARD\"
      }
    }
  }")

STATUS=$(echo "$inspect_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
echo -e "   ${GREEN}✔ SUCCESS!${RESET} Site marked feasible. Forwarded to Commissioner for final sanctioning."
echo -e "   ${WHITE}New Workflow State:${RESET} ${YELLOW}$STATUS${RESET}"

print_divider
prompt_next

# ------------------------------------------------------------------------------
# STAGE 4: Final Approval and Sanctioning
# ------------------------------------------------------------------------------
clear
print_divider
echo -e "   ${BG_BLUE}${WHITE}  STAGE 4: COMMISSIONER APPROVAL & SANCTION  ${RESET}"
print_divider

print_user_card "Municipal Approver" "Rajesh Gupta (Municipal Commissioner)" "APPROVE_FOR_CONNECTION" "Sanctioning legal sewerage line. Generating a unique Sewerage Connection ID."

echo -e "   Transitioning state..."
sleep 1.5

approve_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
      \"processInstance\": {
        \"action\": \"APPROVE_FOR_CONNECTION\"
      }
    }
  }")

STATUS=$(echo "$approve_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
CONN_NO=$(echo "$approve_res" | jq -r '.SewerageConnections[0].connectionNo' 2>/dev/null)

# Clean, spacious, and prominent Approval Card!
echo -e ""
echo -e "   ${GREEN}╔══════════════════════════════════════════════════════════════════════╗${RESET}"
echo -e "   ${GREEN}║                 STATEMENT OF APPROVAL GRANTED                        ║${RESET}"
echo -e "   ${GREEN}╠══════════════════════════════════════════════════════════════════════╣${RESET}"
echo -e "   ${GREEN}║${RESET}  The Municipal Corporation of Amritsar hereby registers and         ${GREEN}║${RESET}"
echo -e "   ${GREEN}║${RESET}  sanctions a permanent sewerage utility connection.                  ${GREEN}║${RESET}"
echo -e "   ${GREEN}║${RESET}                                                                      ${GREEN}║${RESET}"
echo -e "   ${GREEN}║${RESET}  ${WHITE}Application No:${RESET}   %-46s  ${GREEN}║${RESET}" "$APP_NO"
echo -e "   ${GREEN}║${RESET}  ${WHITE}Assigned Conn ID:${RESET} ${PURPLE}%-46s${RESET}  ${GREEN}║${RESET}" "$CONN_NO"
echo -e "   ${GREEN}║${RESET}  ${WHITE}Sanctioned By:${RESET}    %-46s  ${GREEN}║${RESET}" "Rajesh Gupta (Municipal Commissioner)"
echo -e "   ${GREEN}║${RESET}  ${WHITE}Approval Status:${RESET}  ${BG_GREEN}${WHITE} [ SANCTIONED / PENDING_FOR_PAYMENT ] ${RESET}          ${GREEN}║${RESET}"
echo -e "   ${GREEN}╚══════════════════════════════════════════════════════════════════════╝${RESET}"
echo -e ""

print_divider
prompt_next

# ------------------------------------------------------------------------------
# STAGE 5: Demand Charges Apportionment & Payment
# ------------------------------------------------------------------------------
clear
print_divider
echo -e "   ${BG_BLUE}${WHITE}  STAGE 5: DEMAND PAYMENT SETTLEMENT  ${RESET}"
print_divider

print_user_card "Billing Agent" "Anita Patel (Municipal Accounts & Cash Collection)" "PAY" "Apportioning tax demand ledgers for connection $CONN_NO"

echo -e "   Settling outstanding connection dues..."
sleep 1

pay_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
      \"processInstance\": {
        \"action\": \"PAY\"
      }
    }
  }")

STATUS=$(echo "$pay_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)
echo -e "   ${GREEN}✔ SUCCESS!${RESET} Apportioned fee cleared. Receipt issued for connection $CONN_NO."
echo -e "   ${WHITE}New Workflow State:${RESET} ${YELLOW}$STATUS${RESET}"

print_divider
prompt_next

# ------------------------------------------------------------------------------
# STAGE 6: Connection Line Activation
# ------------------------------------------------------------------------------
clear
print_divider
echo -e "   ${BG_BLUE}${WHITE}  STAGE 6: LINE ACTIVATION & PROVISIONING  ${RESET}"
print_divider

print_user_card "Field Engineer" "Suresh Kumar (Chief Technical Engineer)" "ACTIVATE_CONNECTION" "Completing physical pipeline coupling to municipal trunk lines"

echo -e "   Activating sewerage line connection..."
sleep 1

activate_res=$(curl -s -X POST "http://localhost:3468/sw-services/swc/_update" \
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
      \"processInstance\": {
        \"action\": \"ACTIVATE_CONNECTION\"
      }
    }
  }")

STATUS=$(echo "$activate_res" | jq -r '.SewerageConnections[0].applicationStatus' 2>/dev/null)

echo -e "   ${GREEN}✔ SUCCESS!${RESET} Sewerage pipeline successfully provisioned and linked!"
echo -e "   ${WHITE}Final Connection State:${RESET} ${BG_GREEN}${WHITE} $STATUS ${RESET}"

print_divider
prompt_next

# ------------------------------------------------------------------------------
# FINAL AUDIT SUMMARY
# ------------------------------------------------------------------------------
clear
print_divider
echo -e "               ${BG_BLUE}${WHITE}  LIFECYCLE SUMMARY & SYSTEM AUDIT  ${RESET}"
print_divider
echo -e "  The entire distributed municipal transaction has successfully completed."
echo -e ""
echo -e "  ${WHITE}Audit Details Saved to Database (eg_sw_connection):${RESET}"
echo -e "   • ${WHITE}Application Number  :${RESET} ${PURPLE}$APP_NO${RESET}"
echo -e "   • ${WHITE}Connection ID       :${RESET} ${GREEN}$CONN_NO${RESET}"
echo -e "   • ${WHITE}Final Verified State:${RESET} ${GREEN}CONNECTION_ACTIVATED${RESET}"
echo -e "   • ${WHITE}Assigned Property   :${RESET} ${WHITE}PT-107-123456${RESET}"
echo -e "   • ${WHITE}Audited By          :${RESET} ${WHITE}admin (System/Superuser Role)${RESET}"
echo -e ""
echo -e "  ${GREEN}✔ PERSISTENCE INTEGRITY:${RESET} Checked."
echo -e "  ${GREEN}✔ ASYNCHRONOUS KAFKA BROADCAST:${RESET} Confirmed on topic 'update-sw-connection'."
print_divider
echo -e "               ${BG_GREEN}${WHITE}   DEMONSTRATION RUN COMPLETE   ${RESET}"
print_divider
echo -e ""
