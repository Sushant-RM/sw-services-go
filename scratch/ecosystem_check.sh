#!/bin/bash
# ecosystem_check.sh: Structural integrity validation for surrounding DIGIT components

echo -e "\n======================================================================"
echo -e "                 DIGIT ECOSYSTEM STATUS PROBE"
echo -e "======================================================================"

check_port() {
  local host=$1
  local port=$2
  local name=$3
  nc -z -w3 "$host" "$port" 2>/dev/null
  if [ $? -eq 0 ]; then
    echo -e "  * $name ($host:$port) \033[0;32m[ONLINE]\033[0m"
    return 0
  else
    echo -e "  * $name ($host:$port) \033[0;31m[OFFLINE]\033[0m"
    return 1
  fi
}

echo -e "Probing local network container sockets..."
ERRORS=0

check_port "127.0.0.1" 35432 "PostgreSQL Database (sw-postgres)" || ERRORS=$((ERRORS+1))
check_port "127.0.0.1" 39092 "Apache Kafka Broker (sw-kafka)" || ERRORS=$((ERRORS+1))
check_port "127.0.0.1" 3468 "Go Sewerage Service (sw-sw-services)" || ERRORS=$((ERRORS+1))
check_port "127.0.0.1" 3457 "egov-idgen Peer" || ERRORS=$((ERRORS+1))
check_port "127.0.0.1" 3463 "egov-workflow-v2 Peer" || ERRORS=$((ERRORS+1))
check_port "127.0.0.1" 3466 "property-services Peer" || ERRORS=$((ERRORS+1))

echo -e "----------------------------------------------------------------------"
if [ $ERRORS -eq 0 ]; then
  echo -e "ECOSYSTEM STATUS: \033[1;32m[PASS]\033[0m All distributed components ready."
  exit 0
else
  echo -e "ECOSYSTEM STATUS: \033[1;33m[WARNING]\033[0m $ERRORS components are currently unreachable."
  exit 1
fi
