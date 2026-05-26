#!/bin/bash
# db_audit.sh: Performs a physical verification check on the PostgreSQL store

echo -e "\n======================================================================"
echo -e "                 STAGE 4: POSTGRESQL DATABASE AUDIT"
echo -e "======================================================================"

if [ ! -f scratch/.current_app_no ]; then
  echo -e "\033[0;31m[FAIL]\033[0m No active Application Number found in scratch context."
  exit 1
fi

APP_NO=$(cat scratch/.current_app_no)
echo -e "Auditing row persistence for: \033[1;33m$APP_NO\033[0m"

# Execute PostgreSQL query
query_res=$(docker exec -i sw-postgres psql -U postgres -d rainmaker -t -A -c \
  "SELECT applicationno, connectionno, applicationstatus, lastmodifiedby FROM eg_sw_connection WHERE applicationno='$APP_NO';")

if [ -z "$query_res" ]; then
  echo -e "\033[0;31m[FAIL]\033[0m Record not found in PostgreSQL tables!"
  exit 1
fi

IFS='|' read -r app_no conn_no app_status mod_by <<< "$query_res"

echo -e "Physical Schema Column Verification Results:"
echo -e "  * Application Number  : $app_no"
echo -e "  * Connection Number   : $conn_no"
echo -e "  * Application Status   : \033[1;32m$app_status\033[0m"
echo -e "  * Last Modified By    : $mod_by"

if [ "$app_status" == "CONNECTION_ACTIVATED" ]; then
  echo -e "DATABASE INTEGRITY: \033[1;32m[PASS]\033[0m DB record state is perfectly synchronized."
  exit 0
else
  echo -e "DATABASE INTEGRITY: \033[1;31m[FAIL]\033[0m DB record state is out of sync."
  exit 1
fi
